package jwt

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"errors"
)

//Endpoints Wrapper
type Endpoints struct {
	GenerateTokenEndpoint    endpoint.Endpoint
}

//model request
type generateTokenRequest struct {
	Email string `json:"email"`
	Role string `json:"role"`
}

//odel response
type generateTokenResponse struct {
	Token string `json:"token"`
	Err string `json:"err"`
}

//Make actual endpoint per Method
func GenerateTokenEndpoint(svc Service) (endpoint.Endpoint) {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(generateTokenRequest)
		token, err := svc.GenerateToken(ctx, r.Email, r.Role)
		if err != nil {
			return generateTokenResponse{"", err.Error()}, nil
		}
		return generateTokenResponse{token, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) GenerateToken(ctx context.Context, email string, role string) (string, error) {
	resp, err := e.GenerateTokenEndpoint(ctx, generateTokenRequest{Email: email, Role: role})
	if err != nil {
		return "", err
	}
	generateTokenResp := resp.(generateTokenResponse)
	if generateTokenResp.Err != "" {
		return "", errors.New(generateTokenResp.Err)
	}
	return generateTokenResp.Token, nil
}



