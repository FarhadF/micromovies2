package main

import (
	"context"
	"github.com/farhadf/micromovies2/movies/client"
	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"os"
	"strings"
)

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
	moviesService := client.NewGRPCClient(conn)
	if movieId == "" && newMovie == false {
		movies, err := client.GetMovies(ctx, moviesService)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
		//j, err := json.Marshal(mesg)
		logger.Info().Interface("movie", movies).Msg("")
	}
	if movieId != "" && deleteMovie == false {
		movie, err := client.GetMovieById(ctx, moviesService, movieId)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
		logger.Info().Interface("movie", movie).Msg("")

	}
	if newMovie != false && title != "" && director != "" && year != "" && userId != "" {
		dir := strings.Split(director, ",")
		var dirSlice []string
		for _, d := range dir {
			dirSlice = append(dirSlice, d)
		}
		id, err := client.NewMovie(ctx, moviesService, title, dirSlice, year, userId)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
		logger.Info().Str("id", id).Msg("")
	}
	if deleteMovie != false && movieId != "" {
		id, err := client.DeleteMovie(ctx, moviesService, movieId)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		} else {
			logger.Info().Msg("Delete Successful for id: " + id)
		}
	}
	if updateMovie != false && movieId != "" && title != "" && director != "" && year != "" && userId != "" {
		dir := strings.Split(director, ",")
		var dirSlice []string
		for _, d := range dir {
			dirSlice = append(dirSlice, d)
		}
		id, err := client.UpdateMovie(ctx, moviesService, movieId, title, dirSlice, year, userId)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
		logger.Info().Msg("Successfully updated movie with id: " + id)
	}
}
