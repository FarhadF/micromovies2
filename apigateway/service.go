package apigateway

import (
	"google.golang.org/grpc"
	usersClient "micromovies2/users/client"
	"context"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
}

type apigatewayService struct{}

func NewService() Service{
	return apigatewayService{}
}

func (apigatewayService) Login (ctx context.Context, email string, password string) (string, error) {
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

