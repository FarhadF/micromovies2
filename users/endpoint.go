package users

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"errors"
)

//Endpoints Wrapper
type Endpoints struct {
	NewUserEndpoint    endpoint.Endpoint
	GetUserByEmailEndpoint endpoint.Endpoint
}

//model request and response
type newUserResponse struct {
	Id string `json:"id"`
	Err string `json:"err"`
}

//make the actual endpoint
func MakeNewUserEndpoint(svc Service) (endpoint.Endpoint) {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(User)
		id, err := svc.NewUser(ctx, r)
		if err != nil {
			return newUserResponse{"", err.Error()}, nil
		}
		return newUserResponse{id, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) NewUser(ctx context.Context, user User) (string, error) {
	resp, err := e.NewUserEndpoint(ctx, user)
	if err != nil {
		return "", err
	}
	newUserResp := resp.(newUserResponse)
	if newUserResp.Err != "" {
		return "", errors.New(newUserResp.Err)
	}
	return newUserResp.Id, nil
}

//model request and response
type getUserByEmailRequest struct {
	Email string `json:"email"`
}

type getUserByEmailResponse struct {
	User User `json:"user"`
	Err string `json:"err"`
}

//make the actual endpoint
func MakeGetUserByEmailEndpoint(svc Service) (endpoint.Endpoint) {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(getUserByEmailRequest)
		user, err := svc.GetUserByEmail(ctx, r.Email)
		if err != nil {
			return getUserByEmailResponse{user, err.Error()}, nil
		}
		return getUserByEmailResponse{user, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) GetUserByEmail(ctx context.Context, email string) (User, error) {
	req := getUserByEmailRequest{Email: email}
	resp, err := e.GetUserByEmailEndpoint(ctx, req)
	if err != nil {
		return User{}, err
	}
	getUserByEmailResp := resp.(getUserByEmailResponse)
	if getUserByEmailResp.Err != "" {
		return User{}, errors.New(getUserByEmailResp.Err)
	}
	return getUserByEmailResp.User, nil
}