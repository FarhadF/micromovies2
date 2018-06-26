package main

import (
	"context"
	"github.com/farhadf/micromovies2/vault"
	"github.com/farhadf/micromovies2/vault/client"
	"github.com/farhadf/micromovies2/vault/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081", "gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		logger.Fatal().Err(err).Msg("grpc dial err")
	}
	defer conn.Close()
	vaultService := New(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)
	switch cmd {
	case "hash":
		var password string
		password, args = pop(args)
		h, err := client.Hash(ctx, vaultService, password)
		if err != nil {
			logger.Error().Err(err).Msg("")
		}
		logger.Info().Str("hash", h).Msg("")
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		valid, err := client.Validate(ctx, vaultService, password, hash)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if !valid {
			logger.Info().Msg("invalid")
			os.Exit(1)
		}
		logger.Info().Msg("valid")
	default:
		logger.Fatal().Str("unknown command", cmd).Msg("")
	}
}

func New(conn *grpc.ClientConn) vault.Service {
	var hashEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Hash",
		vault.EncodeGRPCHashRequest,
		vault.DecodeGRPCHashResponse,
		pb.HashResponse{},
	).Endpoint()
	var validateEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Validate",
		vault.EncodeGRPCValidateRequest,
		vault.DecodeGRPCValidateResponse,
		pb.ValidateResponse{},
	).Endpoint()
	return vault.Endpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}
