package apigateway

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"micromovies2/movies"
	moviesClient "micromovies2/movies/client"
	"micromovies2/users"
	usersClient "micromovies2/users/client"
)

//business logic of this microservice
type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string, firstname string, lastname string) (string, error)
	ChangePassword(ctx context.Context, email string, oldPassword string, newPassword string) (bool, error)
	GetMovieById(ctx context.Context, id string) (movies.Movie, error)
}

//implementation using empty struct
type apigatewayService struct{}

//create service func, will be used in server.go of this microservice
func NewService() Service {
	return apigatewayService{}
}

//implementation of each method of service interface
func (apigatewayService) Login(ctx context.Context, email string, password string) (string, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("Login", opentracing.ChildOf(span.Context()))
		span.SetTag("email", email)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	conn, err := grpc.Dial(":8084", grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	usersService := usersClient.NewGRPCClient(conn)
	token, err := usersClient.Login(ctx, usersService, email, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

//todo make downstream ports flags/envs
//implementation of each method of service interface
func (apigatewayService) Register(ctx context.Context, email string, password string, firstname string, lastname string) (string, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	usersService := usersClient.NewGRPCClient(conn)
	user := users.User{
		Name:     firstname,
		LastName: lastname,
		Email:    email,
		Password: password,
	}
	id, err := usersClient.NewUser(ctx, usersService, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

//todo make it available to admin and current user only
//implementation of each method of service interface
func (apigatewayService) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (bool, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()
	usersService := usersClient.NewGRPCClient(conn)
	success, err := usersClient.ChangePassword(ctx, usersService, email, currentPassword, newPassword)
	if err != nil {
		return false, err
	}
	return success, nil
}

//implementation of each method of service interface
func (apigatewayService) GetMovieById(ctx context.Context, id string) (movies.Movie, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("GetMovieById", opentracing.ChildOf(span.Context()))
		span.SetTag("id", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return movies.Movie{}, err
	}
	defer conn.Close()
	moviesService := moviesClient.NewGRPCClient(conn)
	movie, err := moviesClient.GetMovieById(ctx, moviesService, id)
	if err != nil {
		return movies.Movie{}, err
	}
	return movie, nil
}
