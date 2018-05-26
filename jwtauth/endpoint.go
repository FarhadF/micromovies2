package jwtauth

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"errors"
)

//Endpoints Wrapper
type Endpoints struct {
	GenerateTokenEndpoint    endpoint.Endpoint
	ParseTokenEndpoint 	endpoint.Endpoint
}

//model request
type generateTokenRequest struct {
	Email string `json:"email"`
	Role string `json:"role"`
}

//model response
type generateTokenResponse struct {
	Token string `json:"token"`
	Err string `json:"err"`
}

//Make actual endpoint per Method
func MakeGenerateTokenEndpoint(svc Service) (endpoint.Endpoint) {
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

//model request
type parseTokenRequest struct {
	Token string `json:"token"`
}



//model response
type parseTokenResponse struct {
	Claims Claims `json:"claims"`
	Err string `json:"err"`
}

//Make actual endpoint per Method
func MakeParseTokenEndpoint(svc Service) (endpoint.Endpoint) {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(parseTokenRequest)
		claims, err := svc.ParseToken(ctx, r.Token)

		if err != nil {
			return parseTokenResponse{Claims{}, err.Error()}, nil
		}
		return parseTokenResponse{claims, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) ParseToken(ctx context.Context, myToken string) (Claims, error) {
	resp, err := e.ParseTokenEndpoint(ctx, parseTokenRequest{Token: myToken})
	if err != nil {
		return Claims{}, err
	}
	parseTokenResp := resp.(parseTokenResponse)
	if parseTokenResp.Err != "" {
		return Claims{}, errors.New(parseTokenResp.Err)
	}
	return parseTokenResp.Claims, nil
}