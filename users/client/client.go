package client

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"micromovies2/users/pb"

	"google.golang.org/grpc/metadata"
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
		//this will inject email in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectEmail),
	).Endpoint()
	var loginEndpoint = grpctransport.NewClient(
		conn, "pb.Users", "Login",
		users.EncodeGRPCLoginRequest,
		users.DecodeGRPCLoginResponse,
		pb.LoginResponse{},
	).Endpoint()
	return users.Endpoints{
		NewUserEndpoint:        newUserEndpoint,
		GetUserByEmailEndpoint: getUserByEmailEndpoint,
		ChangePasswordEndpoint: changePasswordEndpoint,
		LoginEndpoint:          loginEndpoint,
	}
}

func NewUser(ctx context.Context, service users.Service, user users.User) (string, error) {
	id, err := service.NewUser(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetUserByEmail(ctx context.Context, service users.Service, email string) (users.User, error) {
	user, err := service.GetUserByEmail(ctx, email)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func ChangePassword(ctx context.Context, service users.Service, email string, currentPassword string, newPassword string) (bool, error) {
	success, err := service.ChangePassword(ctx, email, currentPassword, newPassword)
	if err != nil {
		return success, err
	}
	return success, nil
}

func Login(ctx context.Context, service users.Service, email string, Password string) (string, error) {
	token, err := service.Login(ctx, email, Password)
	if err != nil {
		return token, err
	}
	return token, nil
}

//client before function
func injectEmail(ctx context.Context, md *metadata.MD) context.Context {
	if email, ok := ctx.Value("email").(string); ok {
		(*md)["email"] = append((*md)["email"], email)
	}
	return ctx
}
