package client

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"micromovies2/movies"
	"micromovies2/movies/pb"
)

//create new client returns GetMovies Service
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
func GetMovies(ctx context.Context, service movies.Service, logger zerolog.Logger) {
	mesg, err := service.GetMovies(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	//j, err := json.Marshal(mesg)
	logger.Info().Interface("movie", mesg).Msg("")
}

func GetMovieById(ctx context.Context, id string, service movies.Service, logger zerolog.Logger) {
	mesg, err := service.GetMovieById(ctx, id)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Interface("movie", mesg).Msg("")
}

func NewMovie(ctx context.Context, title string, director []string, year string, userId string, service movies.Service, logger zerolog.Logger) {
	mesg, err := service.NewMovie(ctx, title, director, year, userId)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Str("id", mesg).Msg("")
}

func DeleteMovie(ctx context.Context, id string, service movies.Service, logger zerolog.Logger) {
	err := service.DeleteMovie(ctx, id)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	} else {
		logger.Info().Msg("Delete Successful for id: " + id)
	}
}

func UpdateMovie(ctx context.Context, id string, title string, director []string, year string, userId string, service movies.Service, logger zerolog.Logger) {
	err := service.UpdateMovie(ctx, id, title, director, year, userId)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Msg("Successfully updated movie with id: " + id)
}
