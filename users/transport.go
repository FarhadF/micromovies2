package users

import (
	"context"
	"github.com/farhadf/micromovies2/users/pb"
)

//encode NewUserRequest
func EncodeGRPCNewUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(User)
	return &pb.NewUserRequest{
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}, nil
}

//decode NewUserRequest
func DecodeGRPCNewUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.NewUserRequest)
	return User{
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}, nil
}

// Encode and Decode NewUserResponse
func EncodeGRPCNewUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(newUserResponse)
	return &pb.NewUserResponse{
		Id:  resp.Id,
		Err: resp.Err,
	}, nil
}

func DecodeGRPCNewUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.NewUserResponse)
	return newUserResponse{
		Id:  resp.Id,
		Err: resp.Err,
	}, nil
}

//encode GetUserByEmailRequest
func EncodeGRPCGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(getUserByEmailRequest)
	return &pb.GetUserByEmailRequest{
		Email: req.Email,
	}, nil
}

//decode GetUserByEmailRequest
func DecodeGRPCGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByEmailRequest)
	return getUserByEmailRequest{
		Email: req.Email,
	}, nil
}

// Encode and Decode GetUserByEmailResponse
func EncodeGRPCGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(getUserByEmailResponse)
	u := &pb.User{
		Id:       resp.User.Id,
		Name:     resp.User.Name,
		LastName: resp.User.LastName,
		Email:    resp.User.Email,
		Role:     resp.User.Role,
	}
	return &pb.GetUserByEmailResponse{
		User: u,
		Err:  resp.Err,
	}, nil
}

func DecodeGRPCGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.GetUserByEmailResponse)
	u := User{
		Id:       resp.User.Id,
		Name:     resp.User.Name,
		LastName: resp.User.LastName,
		Email:    resp.User.Email,
		Role:     resp.User.Role,
	}
	return getUserByEmailResponse{
		User: u,
		Err:  resp.Err,
	}, nil
}

//encode ChangePasswordRequest
func EncodeGRPCChangePasswordRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(changePasswordRequest)
	return &pb.ChangePasswordRequest{
		Email:           req.Email,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}, nil
}

//decode ChangePasswordRequest
func DecodeGRPCChangePasswordRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ChangePasswordRequest)
	return changePasswordRequest{
		Email:           req.Email,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}, nil
}

// Encode and Decode ChangePasswordResponse
func EncodeGRPCChangePasswordResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(changePasswordResponse)
	return &pb.ChangePasswordResponse{
		Success: resp.Success,
		Err:     resp.Err,
	}, nil
}

func DecodeGRPCChangePasswordResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.ChangePasswordResponse)
	return changePasswordResponse{
		Success: resp.Success,
		Err:     resp.Err,
	}, nil
}

//encode LoginRequest
func EncodeGRPCLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(loginRequest)
	return &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

//decode LoginRequest
func DecodeGRPCLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return loginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

// Encode and Decode LoginResponse
func EncodeGRPCLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(loginResponse)
	return &pb.LoginResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		Err:          resp.Err,
	}, nil
}

func DecodeGRPCLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.LoginResponse)
	return loginResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		Err:          resp.Err,
	}, nil
}
