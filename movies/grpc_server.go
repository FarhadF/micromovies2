package movies

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/farhadf/micromovies2/movies/pb"
	"google.golang.org/grpc/metadata"
)

//grpcServer Wrapper
type grpcServer struct {
	getMovies    grpctransport.Handler
	getMovieById grpctransport.Handler
	newMovie     grpctransport.Handler
	deleteMovie  grpctransport.Handler
	updateMovie  grpctransport.Handler
}

// implement getMovies server Interface in movies.pb.go
func (s *grpcServer) GetMovies(ctx context.Context, r *pb.Empty) (*pb.GetMoviesResponse, error) {
	_, resp, err := s.getMovies.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetMoviesResponse), nil
}

// implement getMovieById server Interface in movies.pb.go
func (s *grpcServer) GetMovieById(ctx context.Context, r *pb.GetMovieByIdRequest) (*pb.GetMovieByIdResponse, error) {
	_, resp, err := s.getMovieById.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetMovieByIdResponse), nil
}

// implement NewMovie server Interface in movies.pb.go
func (s *grpcServer) NewMovie(ctx context.Context, r *pb.NewMovieRequest) (*pb.NewMovieResponse, error) {
	_, resp, err := s.newMovie.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.NewMovieResponse), nil
}

// implement DeleteMovie server Interface in movies.pb.go
func (s *grpcServer) DeleteMovie(ctx context.Context, r *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	_, resp, err := s.deleteMovie.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteMovieResponse), nil
}

// implement UpdateMovie server Interface in movies.pb.go
func (s *grpcServer) UpdateMovie(ctx context.Context, r *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	_, resp, err := s.updateMovie.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UpdateMovieResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.MoviesServer {
	return &grpcServer{
		getMovies: grpctransport.NewServer(
			endpoint.GetMoviesEndpoint,
			DecodeGRPCGetMoviesRequest,
			EncodeGRPCGetMoviesResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		getMovieById: grpctransport.NewServer(
			endpoint.GetMovieByIdEndpoint,
			DecodeGRPCGetMovieByIdRequest,
			EncodeGRPCGetMovieByIdResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		newMovie: grpctransport.NewServer(
			endpoint.NewMovieEndpoint,
			DecodeGRPCNewMovieRequest,
			EncodeGRPCNewMovieResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		deleteMovie: grpctransport.NewServer(
			endpoint.DeleteMovieEndpoint,
			DecodeGRPCDeleteMovieRequest,
			EncodeGRPCDeleteMovieResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
		updateMovie: grpctransport.NewServer(
			endpoint.UpdateMovieEndpoint,
			DecodeGRPCUpdateMovieRequest,
			EncodeGRPCUpdateMovieResponse,
			//take out the context set in the upstream service
			grpctransport.ServerBefore(getGRPCContext),
		),
	}
}

//server before: this will retreive email and role from grpc metadata from upstream server and put it in the ctx
func getGRPCContext(ctx context.Context, md metadata.MD) context.Context {
	if email, ok := md["email"]; ok {
		email := email[len(email)-1]
		ctx = context.WithValue(ctx, "email", email)
	}
	if role, ok := md["role"]; ok {
		role := role[len(role)-1]
		ctx = context.WithValue(ctx, "role", role)
	}

	if correlationid, ok := md["correlationid"]; ok {
		correlationid := correlationid[len(correlationid)-1]
		ctx = context.WithValue(ctx, "correlationid", correlationid)
	}
	return ctx
}
