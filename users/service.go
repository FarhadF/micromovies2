package users

import (
	"context"
	jwtClient "github.com/farhadf/micromovies2/jwtauth/client"
	vaultClient "github.com/farhadf/micromovies2/vault/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jackc/pgx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

//business logic methods
type Service interface {
	NewUser(ctx context.Context, user User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) (bool, error)
	Login(ctx context.Context, email string, password string) (string, error)
	//todo: edit user
}

//implementation with database and logger
type usersService struct {
	db     *pgx.ConnPool
	logger zap.Logger
	config Config
}

//constructor - we can later add initialization if needed
func NewService(db *pgx.ConnPool, logger zap.Logger, config Config) Service {
	return usersService{
		db:     db,
		logger: logger,
		config: config,
	}
}

//method implementation
func (s usersService) NewUser(ctx context.Context, user User) (string, error) {
	rows, err := s.db.Query("select * from users where email='" + user.Email + "'")
	defer rows.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if !rows.Next() {
		conn, err := grpc.Dial(s.config.VaultAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
		if err != nil {
			return "", errors.WithStack(err)
		}
		defer conn.Close()
		vaultService := vaultClient.New(conn)
		hash, err := vaultClient.Hash(ctx, vaultService, user.Password)
		if err != nil {
			return "", errors.WithStack(err)
		}
		var id string
		user.Role = "user"
		err = s.db.QueryRow("insert into users (name, lastname, email, password, userrole) values($1,$2,$3,$4,$5) returning id",
			user.Name, user.LastName, user.Email, hash, user.Role).Scan(&id)
		if err != nil {
			return "", errors.WithStack(err)
		}
		//return strconv.FormatInt(id, 10), nil
		return id, nil
	} else {

		return "", errors.New("user already exists")
	}

}

//method implementation
func (s usersService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := s.db.QueryRow("select * from users where email='"+email+"'").Scan(&user.Id, &user.Name, &user.LastName,
		&user.Email, &user.Password, &user.Role, &user.CreatedOn, &user.UpdatedOn, &user.LastLogin)
	if err != nil {
		return user, errors.WithStack(err)
	}
	return user, nil
}

//method implementation
func (s usersService) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (bool, error) {
	var currentPasswordHash string
	err := s.db.QueryRow("select password from users where email='" + email + "'").Scan(&currentPasswordHash)
	if err != nil {
		return false, errors.WithStack(err)
	}
	conn, err := grpc.Dial(s.config.VaultAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer conn.Close()
	vaultService := vaultClient.New(conn)
	valid, err := vaultClient.Validate(ctx, vaultService, currentPassword, currentPasswordHash)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if valid != true {
		return false, errors.WithStack(errors.New("email or password incorrect"))
	}
	hash, err := vaultClient.Hash(ctx, vaultService, newPassword)
	if err != nil {
		return false, errors.WithStack(err)
	}
	_, err = s.db.Exec("update users set password=$1 where email=$2", hash, email)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return true, nil
}

//method implementation
func (s usersService) Login(ctx context.Context, email string, password string) (string, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("Login", opentracing.ChildOf(span.Context()))
		span.SetTag("email", email)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.WithStack(err)
	}
	//call vault service + grpc_opentracing interceptor for client
	conn, err := grpc.Dial(s.config.VaultAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer conn.Close()
	vaultService := vaultClient.New(conn)
	valid, err := vaultClient.Validate(ctx, vaultService, password, user.Password)
	if valid != true {
		return "", errors.WithStack(errors.New("email or password incorrect"))
	}
	//call jwt client + grpc_opentracing interceptor for client
	conn1, err := grpc.Dial(s.config.JwtAuthAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer conn1.Close()
	jwtService := jwtClient.NewGRPCClient(conn1)
	token, err := jwtClient.GenerateToken(ctx, jwtService, email, user.Role)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return token, nil
}

//required Configuration to pass down to the service from the flags in cmd/server.go
type Config struct {
	VaultAddr   string
	JwtAuthAddr string
}
