package main

import (
	"context"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"log"
	"micromovies2/vault"
	"micromovies2/vault/client"
	"micromovies2/vault/pb"
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
			log.Fatalln(err.Error())
		}
		fmt.Println(h)
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		valid, err := client.Validate(ctx, vaultService, password, hash)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if !valid {
			fmt.Println("invalid")
			os.Exit(1)
		}
		fmt.Println("valid")
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
