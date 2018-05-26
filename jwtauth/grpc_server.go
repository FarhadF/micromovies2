package jwtauth

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"context"
	"micromovies2/jwtauth/pb"
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
		),
		parseToken: grpctransport.NewServer(
			endpoint.ParseTokenEndpoint,
			DecodeGRPCParseTokenRequest,
			EncodeGRPCParseTokenResponse,
		),
	}
}
