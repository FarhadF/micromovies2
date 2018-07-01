package jwtauth

import (
	"context"
	"github.com/farhadf/micromovies2/jwtauth/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/metadata"
)

//grpcServer Wrapper
type grpcServer struct {
	generateToken grpctransport.Handler
	parseToken    grpctransport.Handler
}

// implement GenerateToken server Interface in jwt.pb.go
func (s *grpcServer) GenerateToken(ctx context.Context, r *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	_, resp, err := s.generateToken.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GenerateTokenResponse), nil
}

// implement ParseToken server Interface in jtw.pb.go
func (s *grpcServer) ParseToken(ctx context.Context, r *pb.ParseTokenRequest) (*pb.ParseTokenResponse, error) {
	_, resp, err := s.parseToken.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ParseTokenResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.JWTServer {
	return &grpcServer{
		generateToken: grpctransport.NewServer(
			endpoint.GenerateTokenEndpoint,
			DecodeGRPCGenerateTokenRequest,
			EncodeGRPCGenerateTokenResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		parseToken: grpctransport.NewServer(
			endpoint.ParseTokenEndpoint,
			DecodeGRPCParseTokenRequest,
			EncodeGRPCParseTokenResponse,
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
