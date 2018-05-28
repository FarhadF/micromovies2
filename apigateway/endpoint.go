package apigateway

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"errors"
)

//Endpoints Wrapper
type Endpoints struct {
	Ctx              context.Context
	LoginEndpoint          endpoint.Endpoint
}

//model request and response
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"err,omitempty"`
}

//make the actual endpoint
func MakeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(loginRequest)
		token, err := svc.Login(ctx, r.Email, r.Password)
		if err != nil {
			return loginResponse{token, err.Error()}, nil
		}
		return loginResponse{token, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) Login(ctx context.Context, email string, Password string) (string, error) {
	req := loginRequest{Email: email, Password: Password}
	resp, err := e.LoginEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	loginResp := resp.(loginResponse)
	if loginResp.Err != "" {
		return "", errors.New(loginResp.Err)
	}
	return loginResp.Token, nil
}