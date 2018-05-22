package users

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/users/pb"
	"context"
)

type grpcServer struct {
	newUser    grpctransport.Handler
}

// implement NewUser server Interface in movies.pb.go
func (s *grpcServer) NewUser(ctx context.Context, r *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	_, resp, err := s.newUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.NewUserResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.UserServer {
	return &grpcServer{
		newUser: grpctransport.NewServer(
			endpoint.NewUserEndpoint,
			DecodeGRPCNewUserRequest,
			EncodeGRPCNewUserResponse,
		),
	}
}
