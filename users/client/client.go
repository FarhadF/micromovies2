package client

import (
	"google.golang.org/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/users/pb"
	"context"

	"micromovies2/users"
)

func New(conn *grpc.ClientConn) users.Service {
	var newUserEndpoint = grpctransport.NewClient(
		conn, "pb.Users", "NewUser",
		users.EncodeGRPCNewUserRequest,
		users.DecodeGRPCNewUserResponse,
		pb.NewUserResponse{},
	).Endpoint()
	return users.Endpoints{
		NewUserEndpoint:     newUserEndpoint,
	}
}

func NewUser(ctx context.Context, service users.Service, user users.User) (string, error){
	id, err := service.NewUser(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}
