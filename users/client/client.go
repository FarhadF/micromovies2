package client

import (
	"google.golang.org/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/users/pb"
	"context"

	"micromovies2/users"
)

func NewGRPCClient(conn *grpc.ClientConn) users.Service {
	var newUserEndpoint = grpctransport.NewClient(
		conn, "pb.Users", "NewUser",
		users.EncodeGRPCNewUserRequest,
		users.DecodeGRPCNewUserResponse,
		pb.NewUserResponse{},
	).Endpoint()
	var getUserByEmailEndpoint = grpctransport.NewClient(
		conn, "pb.Users", "GetUserByEmail",
		users.EncodeGRPCGetUserByEmailRequest,
		users.DecodeGRPCGetUserByEmailResponse,
		pb.GetUserByEmailResponse{},
	).Endpoint()
	var changePasswordEndpoint = grpctransport.NewClient(
		conn, "pb.Users", "ChangePassword",
		users.EncodeGRPCChangePasswordRequest,
		users.DecodeGRPCChangePasswordResponse,
		pb.ChangePasswordResponse{},
	).Endpoint()
	return users.Endpoints{
		NewUserEndpoint:     newUserEndpoint,
		GetUserByEmailEndpoint: getUserByEmailEndpoint,
		ChangePasswordEndpoint: changePasswordEndpoint,
	}
}

func NewUser(ctx context.Context, service users.Service, user users.User) (string, error){
	id, err := service.NewUser(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetUserByEmail(ctx context.Context, service users.Service, email string) (users.User, error){
	user, err := service.GetUserByEmail(ctx, email)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func ChangePassword(ctx context.Context, service users.Service, email string, currentPassword string, newPassword string) (bool, error){
	success, err := service.ChangePassword(ctx, email, currentPassword, newPassword)
	if err != nil {
		return success, err
	}
	return success, nil
}