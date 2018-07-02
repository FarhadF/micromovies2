package client

import (
	"context"
	"github.com/farhadf/micromovies2/movies"
	"github.com/farhadf/micromovies2/movies/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//create new client returns Movies Service
func NewGRPCClient(conn *grpc.ClientConn) movies.Service {
	var getMoviesEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "GetMovies",
		movies.EncodeGRPCGetMoviesRequest,
		movies.DecodeGRPCGetMoviesResponse,
		pb.GetMoviesResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()
	var getMovieByIdEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "GetMovieById",
		movies.EncodeGRPCGetMovieByIdRequest,
		movies.DecodeGRPCGetMovieByIdResponse,
		pb.GetMovieByIdResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()
	var newMovieEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "NewMovie",
		movies.EncodeGRPCNewMovieRequest,
		movies.DecodeGRPCNewMovieResponse,
		pb.NewMovieResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()
	var deleteMovieEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "DeleteMovie",
		movies.EncodeGRPCDeleteMovieRequest,
		movies.DecodeGRPCDeleteMovieResponse,
		pb.DeleteMovieResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
	).Endpoint()

	var updateMovieEndpoint = grpctransport.NewClient(
		conn, "pb.Movies", "UpdateMovie",
		movies.EncodeGRPCUpdateMovieRequest,
		movies.DecodeGRPCUpdateMovieResponse,
		pb.UpdateMovieResponse{},
		//this will inject context fields specified in injectContext func in the metadata to be passed to downstream service
		grpctransport.ClientBefore(injectContext),
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
func GetMovies(ctx context.Context, service movies.Service) ([]movies.Movie, error) {
	movies, err := service.GetMovies(ctx)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func GetMovieById(ctx context.Context, service movies.Service, id string) (movies.Movie, error) {
	movie, err := service.GetMovieById(ctx, id)
	if err != nil {
		return movies.Movie{}, err
	}
	return movie, nil
}

func NewMovie(ctx context.Context, service movies.Service, title string, director []string, year string, userId string) (string, error) {
	id, err := service.NewMovie(ctx, title, director, year, userId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteMovie(ctx context.Context, service movies.Service, id string) (string, error) {
	err := service.DeleteMovie(ctx, id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func UpdateMovie(ctx context.Context, service movies.Service, id string, title string, director []string, year string, userId string) (string, error) {
	err := service.UpdateMovie(ctx, id, title, director, year, userId)
	if err != nil {
		return "", err
	}
	return id, nil
}

//client before function to inject context into grpc metadata to pass to downstream service
func injectContext(ctx context.Context, md *metadata.MD) context.Context {
	if email, ok := ctx.Value("email").(string); ok {
		(*md)["email"] = append((*md)["email"], email)
	}
	if role, ok := ctx.Value("role").(string); ok {
		(*md)["role"] = append((*md)["role"], role)
	}
	if correlationid, ok := ctx.Value("correlationid").(string); ok {
		(*md)["correlationid"] = append((*md)["correlationid"], correlationid)
	}
	return ctx
}
