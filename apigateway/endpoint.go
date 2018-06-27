package apigateway

import (
	"context"
	"github.com/farhadf/micromovies2/movies"
	"github.com/go-kit/kit/endpoint"
)

//Endpoints Wrapper
type Endpoints struct {
	Ctx                    context.Context
	LoginEndpoint          endpoint.Endpoint
	RegisterEndpoint       endpoint.Endpoint
	ChangePasswordEndpoint endpoint.Endpoint
	GetMovieByIdEndpoint   endpoint.Endpoint
	NewMovieEndpoint       endpoint.Endpoint
}

//model request and response
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// If successful Token will be in response, If it fails, Check err in the response.
// swagger:response loginResponse
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

//model request
type registerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

//model response
type registerResponse struct {
	Id  string `json:"id"`
	Err string `json"err"`
}

//make the actual endpoint
func MakeRegisterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(registerRequest)
		id, err := svc.Register(ctx, r.Email, r.Password, r.FirstName, r.LastName)
		if err != nil {
			return registerResponse{id, err.Error()}, nil
		}
		return registerResponse{id, ""}, nil
	}
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

//model request and response
type getMovieByIdRequest struct {
	Id string `json:"id"`
}

type getMovieByIdResponse struct {
	Movie movies.Movie `json:"movie"`
	Err   string       `json:"err"`
}

//make the actual endpoint
func MakeGetMovieByIdEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		id := req.(string)
		movie, err := svc.GetMovieById(ctx, id)
		if err != nil {
			return getMovieByIdResponse{movies.Movie{}, err.Error()}, nil
		}
		return getMovieByIdResponse{movie, ""}, nil
	}
}

//model request and response
type newMovieRequest struct {
	Title     string   `json:"title"`
	Director  []string `json:"director"`
	Year      string   `json:"year"`
	Createdby string   `json:"createdby"`
}

type newMovieResponse struct {
	Id  string `json:"id"`
	Err string `json:"err"`
}

//make the actual endpoint
func MakeNewMovieEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(newMovieRequest)
		id, err := svc.NewMovie(ctx, r.Title, r.Director, r.Year, r.Createdby)
		if err != nil {
			return newMovieResponse{"", err.Error()}, nil
		}
		return newMovieResponse{id, ""}, nil
	}
}
