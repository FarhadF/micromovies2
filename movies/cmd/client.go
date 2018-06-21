package main

import (
	"context"
	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"micromovies2/movies/client"
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
		client.GetMovies(ctx, moviesService, logger)
	}
	if movieId != "" && deleteMovie == false {
		client.GetMovieById(ctx, movieId, moviesService, logger)

	}
	if newMovie != false && title != "" && director != "" && year != "" && userId != "" {
		dir := strings.Split(director, ",")
		var dirSlice []string
		for _, d := range dir {
			dirSlice = append(dirSlice, d)
		}
		client.NewMovie(ctx, title, dirSlice, year, userId, moviesService, logger)
	}
	if deleteMovie != false && movieId != "" {
		client.DeleteMovie(ctx, movieId, moviesService, logger)
	}
	if updateMovie != false && movieId != "" && title != "" && director != "" && year != "" && userId != "" {
		dir := strings.Split(director, ",")
		var dirSlice []string
		for _, d := range dir {
			dirSlice = append(dirSlice, d)
		}
		client.UpdateMovie(ctx, movieId, title, dirSlice, year, userId, moviesService, logger)
	}

}
