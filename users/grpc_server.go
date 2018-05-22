package users

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/users/pb"
	"context"
)

type grpcServer struct {
	newUser    grpctransport.Handler
	getUserByEmail grpctransport.Handler
}

// implement NewUser server Interface in movies.pb.go
func (s *grpcServer) NewUser(ctx context.Context, r *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	_, resp, err := s.newUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.NewUserResponse), nil
}

// implement NewUser server Interface in movies.pb.go
func (s *grpcServer) GetUserByEmail(ctx context.Context, r *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	_, resp, err := s.getUserByEmail.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetUserByEmailResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.UsersServer {
	return &grpcServer{
		newUser: grpctransport.NewServer(
			endpoint.NewUserEndpoint,
			DecodeGRPCNewUserRequest,
			EncodeGRPCNewUserResponse,
		),
		getUserByEmail: grpctransport.NewServer(
			endpoint.GetUserByEmailEndpoint,
			DecodeGRPCGetUserByEmailRequest,
			EncodeGRPCGetUserByEmailResponse,
		),
	}
}
