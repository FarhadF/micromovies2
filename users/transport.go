package users

import (
	"context"
	"micromovies2/users/pb"
)

//encode NewUserRequest
func EncodeGRPCNewUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(User)
	return &pb.NewUserRequest{
		Name:    req.Name,
		LastName: req.LastName,
		Email:     req.Email,
		Password:   req.Password,
		Role: req.Role,
	}, nil
}

//decode NewUserRequest
func DecodeGRPCNewUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.NewUserRequest)
	return User{
		Name:    req.Name,
		LastName: req.LastName,
		Email:     req.Email,
		Password:   req.Password,
		Role: req.Role,
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
	req := r.(string)
	return &pb.GetUserByEmailRequest{
		Email:     req,
	}, nil
}

//decode GetUserByEmailRequest
func DecodeGRPCGetUserByEmailRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByEmailRequest)
	return User{
		Email:     req.Email,
	}, nil
}

// Encode and Decode GetUserByEmailResponse
func EncodeGRPCGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(getUserByEmailResponse)
	u := &pb.User{
		Id:        resp.User.Id,
		Name:    resp.User.Name,
		LastName: resp.User.LastName,
		Email:     resp.User.Email,
		Role: resp.User.Role,
	}
	return &pb.GetUserByEmailResponse{
		User:  u,
		Err: resp.Err,
	}, nil
}

func DecodeGRPCGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.GetUserByEmailResponse)
	u := User{
		Id:        resp.User.Id,
		Name:    resp.User.Name,
		LastName: resp.User.LastName,
		Email:     resp.User.Email,
		Role: resp.User.Role,
	}
	return getUserByEmailResponse{
		User:  u,
		Err: resp.Err,
	}, nil
}