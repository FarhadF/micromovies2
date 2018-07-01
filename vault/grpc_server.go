package vault

import (
	"context"
	"github.com/farhadf/micromovies2/vault/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	hash     grpctransport.Handler
	validate grpctransport.Handler
}

func (s *grpcServer) Hash(ctx context.Context, r *pb.HashRequest) (*pb.HashResponse, error) {
	_, resp, err := s.hash.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.HashResponse), nil
}

func (s *grpcServer) Validate(ctx context.Context, r *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ValidateResponse), nil
}

func NewGRPCServer(ctx context.Context, endpoints Endpoints) pb.VaultServer {
	return &grpcServer{
		hash: grpctransport.NewServer(
			endpoints.HashEndpoint,
			DecodeGRPCHashRequest,
			EncodeGRPCHashResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		validate: grpctransport.NewServer(
			endpoints.ValidateEndpoint,
			DecodeGRPCValidateRequest,
			EncodeGRPCValidateResponse,
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
