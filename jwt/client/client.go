package client

import (
	"google.golang.org/grpc"
	"micromovies2/jwt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/jwt/pb"
	"context"
)

func NewGRPCClient(conn *grpc.ClientConn) jwt.Service {
	var generateTokenEndpoint = grpctransport.NewClient(
		conn, "pb.JWT", "GenerateToken",
		jwt.EncodeGRPCGenerateTokenRequest,
		jwt.DecodeGRPCGenerateTokenResponse,
		pb.GenerateTokenResponse{},
	).Endpoint()
	return jwt.Endpoints{
		GenerateTokenEndpoint:     generateTokenEndpoint,
	}
}

func GenerateToken(ctx context.Context, service jwt.Service, email string, role string) (string, error){
	h, err := service.GenerateToken(ctx, email, role)
	if err != nil {
		return "", err
	}
	return h, nil
}

