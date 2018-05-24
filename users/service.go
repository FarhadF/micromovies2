package users

import (
	"github.com/jackc/pgx"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"time"
	"micromovies2/vault/client"
	"context"
	"errors"
)

type Service interface{
	NewUser(ctx context.Context, user User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) (bool, error)
	//todo: login
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
func (s usersService) NewUser (ctx context.Context, user User) (string, error) {
	rows, err := s.db.Query("select * from users where email='" + user.Email + "'")
	defer rows.Close()
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
		if err != nil {
			s.logger.Error().Err(err).Msg("grpc dial err")
			return "", err
		}
		defer conn.Close()
		vaultService := client.New(conn)
		hash, err := client.Hash(ctx, vaultService, user.Password)
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
func (s usersService) GetUserByEmail (ctx context.Context, email string) (User, error) {
	var user User
	err := s.db.QueryRow("select * from users where email='" + email + "'").Scan(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

//method implementation
func (s usersService) ChangePassword (ctx context.Context, email string, currentPassword string, newPassword string) (bool, error) {
	var currentPasswordHash string
	err := s.db.QueryRow("select password from users where email='" + email + "'").Scan(&currentPasswordHash)
	if err != nil {
		return false, err
	}
	conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		s.logger.Error().Err(err).Msg("grpc dial err")
		return false, err
	}
	defer conn.Close()
	vaultService := client.New(conn)
	valid, err := client.Validate(ctx, vaultService, currentPassword, currentPasswordHash )
	if err != nil {
		return false, err
	}
	if valid != true {
		return false, errors.New("wrong password")
	}
	hash, err := client.Hash(ctx, vaultService, newPassword)
	if err != nil {
		return false, err
	}
	_, err = s.db.Exec("update users set password=$1 where email=$2", hash, email)
	if err != nil {
		return false, err
	}
	return true, nil
}