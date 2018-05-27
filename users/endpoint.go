package users

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

//Endpoints Wrapper
type Endpoints struct {
	NewUserEndpoint        endpoint.Endpoint
	GetUserByEmailEndpoint endpoint.Endpoint
	ChangePasswordEndpoint endpoint.Endpoint
	LoginEndpoint          endpoint.Endpoint
}

//model request and response
type newUserResponse struct {
	Id  string `json:"id"`
	Err string `json:"err"`
}

//make the actual endpoint
func MakeNewUserEndpoint(svc Service) endpoint.Endpoint {
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
	User User   `json:"user"`
	Err  string `json:"err"`
}

//make the actual endpoint
func MakeGetUserByEmailEndpoint(svc Service) endpoint.Endpoint {
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

//model request and response
type changePasswordRequest struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"currentpassword"`
	NewPassword     string `json:"newpassword"`
}

type changePasswordResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
}

//make the actual endpoint
func MakeChangePasswordEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(changePasswordRequest)
		success, err := svc.ChangePassword(ctx, r.Email, r.CurrentPassword, r.NewPassword)
		if err != nil {
			return changePasswordResponse{success, err.Error()}, nil
		}
		return changePasswordResponse{success, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (bool, error) {
	req := changePasswordRequest{Email: email, CurrentPassword: currentPassword, NewPassword: newPassword}
	resp, err := e.ChangePasswordEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	changePasswordResp := resp.(changePasswordResponse)
	if changePasswordResp.Err != "" {
		return false, errors.New(changePasswordResp.Err)
	}
	return changePasswordResp.Success, nil
}

//model request and response
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
	Err   string `json:"err"`
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
