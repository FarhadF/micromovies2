package users

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/users/pb"
	"google.golang.org/grpc/metadata"
	"fmt"
)

type grpcServer struct {
	newUser        grpctransport.Handler
	getUserByEmail grpctransport.Handler
	changePassword grpctransport.Handler
	login          grpctransport.Handler
}

// implement NewUser server Interface in movies.pb.go
func (s *grpcServer) NewUser(ctx context.Context, r *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	_, resp, err := s.newUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.NewUserResponse), nil
}

// implement GetUserByEmail server Interface in users.pb.go
func (s *grpcServer) GetUserByEmail(ctx context.Context, r *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	_, resp, err := s.getUserByEmail.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetUserByEmailResponse), nil
}

// implement ChangePassword server Interface in users.pb.go
func (s *grpcServer) ChangePassword(ctx context.Context, r *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	_, resp, err := s.changePassword.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ChangePasswordResponse), nil
}

// implement Login server Interface in users.pb.go
func (s *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.UsersServer {
	return &grpcServer{
		newUser: grpctransport.NewServer(
			endpoint.NewUserEndpoint,
			DecodeGRPCNewUserRequest,
			EncodeGRPCNewUserResponse,
			grpctransport.ServerBefore(getGRPCContext),
		),
		getUserByEmail: grpctransport.NewServer(
			endpoint.GetUserByEmailEndpoint,
			DecodeGRPCGetUserByEmailRequest,
			EncodeGRPCGetUserByEmailResponse,
		),
		changePassword: grpctransport.NewServer(
			endpoint.ChangePasswordEndpoint,
			DecodeGRPCChangePasswordRequest,
			EncodeGRPCChangePasswordResponse,
		),
		login: grpctransport.NewServer(
			endpoint.LoginEndpoint,
			DecodeGRPCLoginRequest,
			EncodeGRPCLoginResponse,
		),
	}
}

//serverbefore funcs
func getGRPCContext(ctx context.Context, md metadata.MD) context.Context {
	if hdr, ok := md[string("email")]; ok {
		email := hdr[len(hdr)-1]
		ctx = context.WithValue(ctx, "email", email)
		fmt.Printf("\tServer received email %q in metadata header, set context\n", email)
	}
	return ctx
}
