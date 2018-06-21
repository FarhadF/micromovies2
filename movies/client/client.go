package client

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"micromovies2/movies"
	"micromovies2/movies/pb"
)

//create new client returns Movies Service
func NewGRPCClient(conn *grpc.ClientConn) movies.Service {
	var getMoviesEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "GetMovies",
		movies.EncodeGRPCGetMoviesRequest,
		movies.DecodeGRPCGetMoviesResponse,
		pb.GetMoviesResponse{},
	).Endpoint()
	var getMovieByIdEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "GetMovieById",
		movies.EncodeGRPCGetMovieByIdRequest,
		movies.DecodeGRPCGetMovieByIdResponse,
		pb.GetMovieByIdResponse{},
	).Endpoint()
	var newMovieEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "NewMovie",
		movies.EncodeGRPCNewMovieRequest,
		movies.DecodeGRPCNewMovieResponse,
		pb.NewMovieResponse{},
	).Endpoint()
	var deleteMovieEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "DeleteMovie",
		movies.EncodeGRPCDeleteMovieRequest,
		movies.DecodeGRPCDeleteMovieResponse,
		pb.DeleteMovieResponse{},
	).Endpoint()
	var updateMovieEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "UpdateMovie",
		movies.EncodeGRPCUpdateMovieRequest,
		movies.DecodeGRPCUpdateMovieResponse,
		pb.UpdateMovieResponse{},
	).Endpoint()
	return movies.Endpoints{
		GetMoviesEndpoint:    getMoviesEndpoint,
		GetMovieByIdEndpoint: getMovieByIdEndpoint,
		NewMovieEndpoint:     newMovieEndpoint,
		DeleteMovieEndpoint:  deleteMovieEndpoint,
		UpdateMovieEndpoint:  updateMovieEndpoint,
	}
}

//callService helper
func GetMovies(ctx context.Context, service movies.Service)([]movies.Movie, error) {
	movies, err := service.GetMovies(ctx)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func GetMovieById(ctx context.Context, id string, service movies.Service)(movies.Movie, error) {
	movie, err := service.GetMovieById(ctx, id)
	if err != nil {
		return movies.Movie{}, err
	}
	return movie, nil
}

func NewMovie(ctx context.Context, title string, director []string, year string, userId string, service movies.Service)(string, error) {
	id, err := service.NewMovie(ctx, title, director, year, userId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteMovie(ctx context.Context, id string, service movies.Service)(string, error){
	err := service.DeleteMovie(ctx, id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func UpdateMovie(ctx context.Context, id string, title string, director []string, year string, userId string, service movies.Service)(string, error) {
	err := service.UpdateMovie(ctx, id, title, director, year, userId)
	if err != nil {
		return "", err
	}
	return id, nil
}