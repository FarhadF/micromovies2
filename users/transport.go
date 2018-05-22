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

