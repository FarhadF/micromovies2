package movies

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

//model request and response
type getMoviesResponse struct {
	Movies []Movie `json:"movies,omitempty"`
	Err    string  `json:"err,omitempty"`
}

//Endpoints Wrapper
type Endpoints struct {
	GetMoviesEndpoint    endpoint.Endpoint
	GetMovieByIdEndpoint endpoint.Endpoint
	NewMovieEndpoint     endpoint.Endpoint
	DeleteMovieEndpoint  endpoint.Endpoint
	UpdateMovieEndpoint  endpoint.Endpoint
}

//Make actual endpoint per Method
func MakeGetMoviesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		movies, err := svc.GetMovies(ctx)
		if err != nil {
			return getMoviesResponse{nil, err.Error()}, nil
		}
		return getMoviesResponse{movies, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) GetMovies(ctx context.Context) ([]Movie, error) {

	resp, err := e.GetMoviesEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	getMoviesResp := resp.(getMoviesResponse)
	if getMoviesResp.Err != "" {
		return nil, errors.New(getMoviesResp.Err)
	}
	return getMoviesResp.Movies, nil
}

//model request and response
type getMovieByIdRequest struct {
	Id string `json:"id"`
}

type getMovieByIdResponse struct {
	Movie Movie  `json:"movie"`
	Err   string `json:"err"`
}

//Make actual endpoint per Method
func MakeGetMovieByIdEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(getMovieByIdRequest)
		movie, err := svc.GetMovieById(ctx, r.Id)
		if err != nil {
			return getMovieByIdResponse{movie, err.Error()}, nil
		}
		return getMovieByIdResponse{movie, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) GetMovieById(ctx context.Context, id string) (Movie, error) {
	req := getMovieByIdRequest{
		Id: id,
	}
	var movie Movie
	resp, err := e.GetMovieByIdEndpoint(ctx, req)
	if err != nil {
		return movie, err
	}
	getMovieByIdResp := resp.(getMovieByIdResponse)
	if getMovieByIdResp.Err != "" {
		return movie, errors.New(getMovieByIdResp.Err)
	}
	return getMovieByIdResp.Movie, nil
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

//Make actual endpoint per Method
func MakeNewMovieEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(newMovieRequest)
		id, err := svc.NewMovie(ctx, r.Title, r.Director, r.Year, r.Createdby)
		if err != nil {
			return newMovieResponse{id, err.Error()}, nil
		}
		return newMovieResponse{id, ""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) NewMovie(ctx context.Context, title string, director []string, year string, createdBy string) (string, error) {
	req := newMovieRequest{
		Title:     title,
		Director:  director,
		Year:      year,
		Createdby: createdBy,
	}
	resp, err := e.NewMovieEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	newMovieResp := resp.(newMovieResponse)
	if newMovieResp.Err != "" {
		return newMovieResp.Id, errors.New(newMovieResp.Err)
	}
	return newMovieResp.Id, nil
}

//model request and response
type deleteMovieRequest struct {
	Id string `json:"id"`
}

type deleteMovieResponse struct {
	Err string `json:"err"`
}

//Make actual endpoint per Method
func MakeDeleteMovieEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(deleteMovieRequest)
		err := svc.DeleteMovie(ctx, r.Id)
		if err != nil {
			return deleteMovieResponse{err.Error()}, nil
		}
		return deleteMovieResponse{""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) DeleteMovie(ctx context.Context, id string) error {
	req := deleteMovieRequest{
		Id: id,
	}
	resp, err := e.DeleteMovieEndpoint(ctx, req)
	if err != nil {
		return err
	}
	deleteMovieResp := resp.(deleteMovieResponse)
	if deleteMovieResp.Err != "" {
		return errors.New(deleteMovieResp.Err)
	}
	return nil
}

//model request and response
type updateMovieRequest struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Director  []string `json:"director"`
	Year      string   `json:"year"`
	Createdby string   `json:"createdby"`
}

type updateMovieResponse struct {
	Err string `json:"err"`
}

//Make actual endpoint per Method
func MakeUpdateMovieEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(updateMovieRequest)
		err := svc.UpdateMovie(ctx, r.Id, r.Title, r.Director, r.Year, r.Createdby)
		if err != nil {
			return updateMovieResponse{err.Error()}, nil
		}
		return updateMovieResponse{""}, nil
	}
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) UpdateMovie(ctx context.Context, id string, title string, director []string, year string, createdBy string) error {
	req := updateMovieRequest{
		Id:        id,
		Title:     title,
		Director:  director,
		Year:      year,
		Createdby: createdBy,
	}
	resp, err := e.UpdateMovieEndpoint(ctx, req)
	if err != nil {
		return err
	}
	updateMovieResp := resp.(updateMovieResponse)
	if updateMovieResp.Err != "" {
		return errors.New(updateMovieResp.Err)
	}
	return nil
}
