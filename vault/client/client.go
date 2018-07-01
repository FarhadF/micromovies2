package client

import (
	"context"
	"github.com/farhadf/micromovies2/vault"
	"github.com/farhadf/micromovies2/vault/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func New(conn *grpc.ClientConn) vault.Service {
	var hashEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Hash",
		vault.EncodeGRPCHashRequest,
		vault.DecodeGRPCHashResponse,
		pb.HashResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()
	var validateEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Validate",
		vault.EncodeGRPCValidateRequest,
		vault.DecodeGRPCValidateResponse,
		pb.ValidateResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
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

//client before function to inject context into grpc metadata to pass to downstream service
func injectContext(ctx context.Context, md *metadata.MD) context.Context {
	if email, ok := ctx.Value("email").(string); ok {
		(*md)["email"] = append((*md)["email"], email)
	}
	if role, ok := ctx.Value("role").(string); ok {
		(*md)["role"] = append((*md)["role"], role)
	}
	if correlationid, ok := ctx.Value("correlationid").(string); ok {
		(*md)["correlationid"] = append((*md)["correlationid"], correlationid)
	}
	return ctx
}
