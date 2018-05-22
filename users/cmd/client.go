package main

import (
	"google.golang.org/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"micromovies2/users/pb"
	flag "github.com/spf13/pflag"
	"context"
	"github.com/rs/zerolog"
	"os"
	"micromovies2/users"
	"micromovies2/users/client"
)

//create new client returns GetMovies Service
func NewGRPCClient(conn *grpc.ClientConn) users.Service {
	var newUserEndpoint = grpctransport.NewClient(
		conn, "pb.User", "NewUser",
		users.EncodeGRPCNewUserRequest,
		users.DecodeGRPCNewUserResponse,
		pb.NewUserResponse{},
	).Endpoint()

	return users.Endpoints{
		NewUserEndpoint: newUserEndpoint,
		}
}

func main() {
	var (
		grpcAddr string
		newUser bool
		name    string
		lastname string
		email     string
		password   string
		role	string
	)
	flag.StringVarP(&grpcAddr, "addr", "a", ":8084", "gRPC address")
	flag.StringVarP(&name, "name", "f", "", "name")
	flag.StringVarP(&lastname, "lastname", "l", "", "lastname")
	flag.StringVarP(&email, "email", "e", "", "email")
	flag.StringVarP(&password, "password", "p", "", "password")
	flag.StringVarP(&role, "role", "r", "user", "role")
	flag.BoolVarP(&newUser, "newuser", "n", false, "newUser")
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
	usersService := NewGRPCClient(conn)
	if newUser == true && name != "" && lastname != "" && email != "" && password != ""{
		user := users.User{Name: name, LastName: lastname, Email: email, Password: password}
		id, err := client.NewUser(ctx, usersService, user)
		if err != nil {
			logger.Error().Err(err).Msg("")
		}
		logger.Info().Msg(id)
	}
	if newUser == false && name == "" && lastname == "" && email != "" && password == "" {
		user, err := client.GetUserByEmail(ctx, usersService, email)
		if err != nil {
			logger.Error().Err(err).Msg("")
		}
		logger.Info().Interface("user",user).Msg("")
	}
}

