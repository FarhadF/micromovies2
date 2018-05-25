package client

import (
	"google.golang.org/grpc"
	"micromovies2/jwtauth"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/jwtauth/pb"
	"context"
)

func NewGRPCClient(conn *grpc.ClientConn) jwtauth.Service {
	var generateTokenEndpoint = grpctransport.NewClient(
		conn, "pb.JWT", "GenerateToken",
		jwtauth.EncodeGRPCGenerateTokenRequest,
		jwtauth.DecodeGRPCGenerateTokenResponse,
		pb.GenerateTokenResponse{},
	).Endpoint()
	return jwtauth.Endpoints{
		GenerateTokenEndpoint:     generateTokenEndpoint,
	}
}

func GenerateToken(ctx context.Context, service jwtauth.Service, email string, role string) (string, error){
	h, err := service.GenerateToken(ctx, email, role)
	if err != nil {
		return "", err
	}
	return h, nil
}

