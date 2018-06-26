package client

import (
	"context"
	"github.com/farhadf/micromovies2/vault"
	"github.com/farhadf/micromovies2/vault/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) vault.Service {
	var hashEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Hash",
		vault.EncodeGRPCHashRequest,
		vault.DecodeGRPCHashResponse,
		pb.HashResponse{},
	).Endpoint()
	var validateEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Validate",
		vault.EncodeGRPCValidateRequest,
		vault.DecodeGRPCValidateResponse,
		pb.ValidateResponse{},
	).Endpoint()
	return vault.Endpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
}

func Hash(ctx context.Context, service vault.Service, password string) (string, error) {
	h, err := service.Hash(ctx, password)
	if err != nil {
		return "", err
	}
	return h, nil
}

func Validate(ctx context.Context, service vault.Service, password, hash string) (bool, error) {
	valid, err := service.Validate(ctx, password, hash)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, nil
	}
	return true, nil
}
