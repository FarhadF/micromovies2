package main

import (
	"google.golang.org/grpc"
	flag "github.com/spf13/pflag"
	"context"
	"github.com/rs/zerolog"
	"os"
	"micromovies2/users"
	"micromovies2/users/client"
)

func main() {
	var (
		grpcAddr       string
		newUser        bool
		name           string
		lastname       string
		email          string
		password       string
		role           string
		changePassword bool
		newPassword    string
	)
	flag.StringVarP(&grpcAddr, "addr", "a", ":8084", "gRPC address")
	flag.StringVarP(&name, "name", "f", "", "name")
	flag.StringVarP(&lastname, "lastname", "l", "", "lastname")
	flag.StringVarP(&email, "email", "e", "", "email")
	flag.StringVarP(&password, "password", "p", "", "password")
	flag.StringVarP(&role, "role", "r", "user", "role")
	flag.BoolVarP(&newUser, "newuser", "n", false, "newUser")
	flag.BoolVarP(&changePassword, "changepassword", "c", false, "changePassword")
	flag.StringVarP(&newPassword, "newpassword", "b", "", "newPassword")
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
	usersService := client.NewGRPCClient(conn)
	if newUser == true && name != "" && lastname != "" && email != "" && password != "" {
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
		} else {
			logger.Info().Interface("user", user).Msg("")
		}
	}
	if changePassword == true && email != "" && password != "" && newPassword != "" {
		success, err := client.ChangePassword(ctx, usersService, email, password, newPassword)
		if err != nil {
			logger.Error().Err(err).Msg("")
		} else {
			logger.Info().Interface("success", success).Msg("")
		}
	}
}
