package users

import (
	"context"
	"github.com/farhadf/micromovies2/users/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/metadata"
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
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		login: grpctransport.NewServer(
			endpoint.LoginEndpoint,
			DecodeGRPCLoginRequest,
			EncodeGRPCLoginResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
	}
}

//server before: this will retreive email and role from grpc metadata from upstream server and put it in the ctx
func getGRPCContext(ctx context.Context, md metadata.MD) context.Context {
	if email, ok := md["email"]; ok {
		email := email[len(email)-1]
		ctx = context.WithValue(ctx, "email", email)
	}
	if role, ok := md["role"]; ok {
		role := role[len(role)-1]
		ctx = context.WithValue(ctx, "role", role)
	}

	if correlationid, ok := md["correlationid"]; ok {
		correlationid := correlationid[len(correlationid)-1]
		ctx = context.WithValue(ctx, "correlationid", correlationid)
	}
	return ctx
}
