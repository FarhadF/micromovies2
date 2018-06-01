package apigateway

import (
	"context"
	"google.golang.org/grpc"
	"micromovies2/users"
	usersClient "micromovies2/users/client"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string, firstname string, lastname string) (string, error)
	ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) (bool, error)
}

type apigatewayService struct{}

func NewService() Service {
	return apigatewayService{}
}

func (apigatewayService) Login(ctx context.Context, email string, password string) (string, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	usersService := usersClient.NewGRPCClient(conn)
	token, err := usersClient.Login(ctx, usersService, email, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (apigatewayService) Register(ctx context.Context, email string, password string, firstname string, lastname string) (string, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	usersService := usersClient.NewGRPCClient(conn)
	user := users.User{
		Name:     firstname,
		LastName: lastname,
		Email:    email,
		Password: password,
	}
	id, err := usersClient.NewUser(ctx, usersService, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (apigatewayService) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (bool, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()
	usersService := usersClient.NewGRPCClient(conn)
	success, err := usersClient.ChangePassword(ctx, usersService, email, currentPassword, newPassword)
	if err != nil {
		return false, err
	}
	return success, nil
}
