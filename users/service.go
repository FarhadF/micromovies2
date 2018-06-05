package users

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	jwtClient "micromovies2/jwtauth/client"
	vaultClient "micromovies2/vault/client"
	"time"
	"google.golang.org/grpc/metadata"
	"fmt"
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
}

//constructor - we can later add initialization if needed
func NewService(db *pgx.ConnPool, logger zap.Logger) Service {
	return usersService{
		db,
		logger,
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
		conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
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
	conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
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
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.WithStack(err)
	}
	conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer conn.Close()
	vaultService := vaultClient.New(conn)
	valid, err := vaultClient.Validate(ctx, vaultService, password, user.Password)
	if valid != true {
		return "", errors.WithStack(errors.New("email or password incorrect"))
	}
	conn1, err := grpc.Dial(":8087", grpc.WithInsecure())
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

func injectEmail(ctx context.Context, md *metadata.MD) context.Context {
	if hdr, ok := ctx.Value(correlationID).(string); ok {
		fmt.Printf("\tClient found correlationID %q in context, set metadata header\n", hdr)
		(*md)[string(correlationID)] = append((*md)[string(correlationID)], hdr)
	}
	return ctx
}