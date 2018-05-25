package jwt

import (
	"context"
	"micromovies2/jwt/pb"
)

//encode GenerateTokenRequest
func EncodeGRPCGenerateTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(generateTokenRequest)
	return &pb.GenerateTokenRequest{
		Email: req.Email,
		Role: req.Role,
	}, nil
}

//decode GenerateTokenRequest
func DecodeGRPCGenerateTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GenerateTokenRequest)
	return generateTokenRequest{
		Email: req.Email,
		Role: req.Role,
	}, nil
}

// encode GenerateToken Response
func EncodeGRPCGenerateTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(generateTokenResponse)
	return &pb.GenerateTokenResponse{
		Token: resp.Token,
		Err:   resp.Err,
	}, nil
}

// decode GenerateToken Response
func DecodeGRPCGenerateTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.GenerateTokenResponse)
	return generateTokenResponse{
		Token: resp.Token,
		Err:   resp.Err,
	}, nil
}
