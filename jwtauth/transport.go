package jwtauth

import (
	"context"
	"github.com/farhadf/micromovies2/jwtauth/pb"
)

//encode GenerateTokenRequest
func EncodeGRPCGenerateTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(generateTokenRequest)
	return &pb.GenerateTokenRequest{
		Email: req.Email,
		Role:  req.Role,
	}, nil
}

//decode GenerateTokenRequest
func DecodeGRPCGenerateTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GenerateTokenRequest)
	return generateTokenRequest{
		Email: req.Email,
		Role:  req.Role,
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

//encode ParseTokenRequest
func EncodeGRPCParseTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(parseTokenRequest)
	return &pb.ParseTokenRequest{
		Token: req.Token,
	}, nil
}

//decode ParseTokenRequest
func DecodeGRPCParseTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ParseTokenRequest)
	return parseTokenRequest{
		Token: req.Token,
	}, nil
}

// encode ParseToken Response
func EncodeGRPCParseTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(parseTokenResponse)
	claims := &pb.Claims{
		Exp:   resp.Claims.Exp,
		Iat:   resp.Claims.Iat,
		Email: resp.Claims.Email,
		Role:  resp.Claims.Role,
	}
	return &pb.ParseTokenResponse{
		Claims: claims,
		Err:    resp.Err,
	}, nil
}

// decode ParseToken Response
func DecodeGRPCParseTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.ParseTokenResponse)
	claims := Claims{
		Exp:   resp.Claims.Exp,
		Iat:   resp.Claims.Iat,
		Email: resp.Claims.Email,
		Role:  resp.Claims.Role,
	}
	return parseTokenResponse{
		Claims: claims,
		Err:    resp.Err,
	}, nil
}
