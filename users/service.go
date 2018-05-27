package users

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	jwtClient "micromovies2/jwtauth/client"
	vaultClient "micromovies2/vault/client"
	"time"
)

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
	logger zerolog.Logger
}

//constructor - we can later add initialization if needed
func NewService(db *pgx.ConnPool, logger zerolog.Logger) Service {
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
		return "", err
	}
	if !rows.Next() {
		conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
		if err != nil {
			return "", err
		}
		defer conn.Close()
		vaultService := vaultClient.New(conn)
		hash, err := vaultClient.Hash(ctx, vaultService, user.Password)
		if err != nil {
			return "", err
		}
		var id string
		user.Role = "user"
		err = s.db.QueryRow("insert into users (name, lastname, email, password, userrole) values($1,$2,$3,$4,$5) returning id",
			user.Name, user.LastName, user.Email, hash, user.Role).Scan(&id)
		if err != nil {
			return "", err
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
		return user, err
	}
	return user, nil
}

//method implementation
func (s usersService) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (bool, error) {
	var currentPasswordHash string
	err := s.db.QueryRow("select password from users where email='" + email + "'").Scan(&currentPasswordHash)
	if err != nil {
		return false, err
	}
	conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return false, err
	}
	defer conn.Close()
	vaultService := vaultClient.New(conn)
	valid, err := vaultClient.Validate(ctx, vaultService, currentPassword, currentPasswordHash)
	if err != nil {
		return false, err
	}
	if valid != true {
		return false, errors.New("wrong password")
	}
	hash, err := vaultClient.Hash(ctx, vaultService, newPassword)
	if err != nil {
		return false, err
	}
	_, err = s.db.Exec("update users set password=$1 where email=$2", hash, email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s usersService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	vaultService := vaultClient.New(conn)
	valid, err := vaultClient.Validate(ctx, vaultService, password, user.Password)
	if valid != true {
		return "", errors.New("wrong password")
	}
	conn1, err := grpc.Dial(":8087", grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn1.Close()
	jwtService := jwtClient.NewGRPCClient(conn1)
	token, err := jwtClient.GenerateToken(ctx, jwtService, email, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
