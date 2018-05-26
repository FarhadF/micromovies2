package main

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"micromovies2/movies"
	"micromovies2/movies/pb"
	"os"
	"strings"
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

func main() {
	var (
		grpcAddr    string
		movieId     string
		newMovie    bool
		title       string
		director    string
		year        string
		userId      string
		deleteMovie bool
		updateMovie bool
	)
	flag.StringVarP(&grpcAddr, "addr", "a", ":8081", "gRPC address")
	flag.StringVarP(&movieId, "id", "i", "", "movieId")
	flag.StringVarP(&title, "title", "t", "", "title")
	flag.StringVarP(&director, "director", "d", "", "director(s) comma seperated")
	flag.StringVarP(&year, "year", "y", "", "year")
	flag.StringVarP(&userId, "userid", "u", "", "userId")
	flag.BoolVarP(&newMovie, "newmovie", "n", false, "newMovie")
	flag.BoolVarP(&deleteMovie, "deletemovie", "D", false, "deleteMovie")
	flag.BoolVarP(&updateMovie, "updatemovie", "U", false, "updateMovie")
	//flag.StringVarP(&requestType, "requestType", "r", "word", "Should be word, sentence or paragraph")
	//flag.IntVarP(&min,"min", "m", 5, "minimum value")
	//flag.IntVarP(&max,"Max", "M", 10, "Maximum value")

	flag.Parse()
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	ctx := context.Background()
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer conn.Close()
	moviesService := NewGRPCClient(conn)
	if movieId == "" && newMovie == false {
		callGetMovies(ctx, moviesService, logger)
	}
	if movieId != "" && deleteMovie == false {
		callGetMovieById(ctx, movieId, moviesService, logger)

	}
	if newMovie != false && title != "" && director != "" && year != "" && userId != "" {
		dir := strings.Split(director, ",")
		var dirSlice []string
		for _, d := range dir {
			dirSlice = append(dirSlice, d)
		}
		callNewMovie(ctx, title, dirSlice, year, userId, moviesService, logger)
	}
	if deleteMovie != false && movieId != "" {
		callDeleteMovie(ctx, movieId, moviesService, logger)
	}
	if updateMovie != false && movieId != "" && title != "" && director != "" && year != "" && userId != "" {
		dir := strings.Split(director, ",")
		var dirSlice []string
		for _, d := range dir {
			dirSlice = append(dirSlice, d)
		}
		callUpdateMovie(ctx, movieId, title, dirSlice, year, userId, moviesService, logger)
	}

}

//callService helper
func callGetMovies(ctx context.Context, service movies.Service, logger zerolog.Logger) {
	mesg, err := service.GetMovies(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	//j, err := json.Marshal(mesg)
	logger.Info().Interface("movie", mesg).Msg("")
}

func callGetMovieById(ctx context.Context, id string, service movies.Service, logger zerolog.Logger) {
	mesg, err := service.GetMovieById(ctx, id)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Interface("movie", mesg).Msg("")
}

func callNewMovie(ctx context.Context, title string, director []string, year string, userId string, service movies.Service, logger zerolog.Logger) {
	mesg, err := service.NewMovie(ctx, title, director, year, userId)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Str("id", mesg).Msg("")
}

func callDeleteMovie(ctx context.Context, id string, service movies.Service, logger zerolog.Logger) {
	err := service.DeleteMovie(ctx, id)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	} else {
		logger.Info().Msg("Delete Successful for id: " + id)
	}
}

func callUpdateMovie(ctx context.Context, id string, title string, director []string, year string, userId string, service movies.Service, logger zerolog.Logger) {
	err := service.UpdateMovie(ctx, id, title, director, year, userId)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Msg("Successfully updated movie with id: " + id)
}
