package client

import (
	"context"
	"github.com/farhadf/micromovies2/jwtauth"
	"github.com/farhadf/micromovies2/jwtauth/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGRPCClient(conn *grpc.ClientConn) jwtauth.Service {
	var generateTokenEndpoint = grpctransport.NewClient(
		conn, "pb.JWT", "GenerateToken",
		jwtauth.EncodeGRPCGenerateTokenRequest,
		jwtauth.DecodeGRPCGenerateTokenResponse,
		pb.GenerateTokenResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()
	var parseTokenEndpoint = grpctransport.NewClient(
		conn, "pb.JWT", "ParseToken",
		jwtauth.EncodeGRPCParseTokenRequest,
		jwtauth.DecodeGRPCParseTokenResponse,
		pb.ParseTokenResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()
	return jwtauth.Endpoints{
		GenerateTokenEndpoint: generateTokenEndpoint,
		ParseTokenEndpoint:    parseTokenEndpoint,
	}
}

func GenerateToken(ctx context.Context, service jwtauth.Service, email string, role string) (string, error) {
	h, err := service.GenerateToken(ctx, email, role)
	if err != nil {
		return "", err
	}
	return h, nil
}

func ParseToken(ctx context.Context, service jwtauth.Service, token string) (jwtauth.Claims, error) {
	h, err := service.ParseToken(ctx, token)
	if err != nil {
		return jwtauth.Claims{}, err
	}
	return h, nil
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
