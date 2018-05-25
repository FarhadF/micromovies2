package main

import (
	"google.golang.org/grpc"
	flag "github.com/spf13/pflag"
	"context"
	"github.com/rs/zerolog"
	"os"
	"micromovies2/jwt/client"
)

func main() {
	var (
		grpcAddr      string
		email         string
		generateToken bool
		role          string
	)
	flag.StringVarP(&grpcAddr, "addr", "a", ":8087", "gRPC address")
	flag.StringVarP(&email, "email", "e", "", "email")
	flag.StringVarP(&role, "role", "r", "user", "role")
	flag.BoolVarP(&generateToken, "generatetoken", "g", false, "generateToken")

	flag.Parse()
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	ctx := context.Background()
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer conn.Close()
	jwtService := client.NewGRPCClient(conn)
	if email != "" && generateToken != false {
		token, err := client.GenerateToken(ctx, jwtService, email, role)
		if err != nil {
			logger.Error().Err(err)
		} else {
			logger.Info().Str("token", token).Msg("")
		}

	}
}
