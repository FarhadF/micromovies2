package jwt

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"context"
	"micromovies2/jwt/pb"
)
//grpcServer Wrapper
type grpcServer struct {
	generateToken    grpctransport.Handler
}

// implement getMovies server Interface in movies.pb.go
func (s *grpcServer) GenerateToken(ctx context.Context, r *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	_, resp, err := s.generateToken.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GenerateTokenResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.JWTServer {
	return &grpcServer{
		generateToken: grpctransport.NewServer(
			endpoint.GenerateTokenEndpoint,
			DecodeGRPCGenerateTokenRequest,
			EncodeGRPCGenerateTokenResponse,
		),
	}
}

