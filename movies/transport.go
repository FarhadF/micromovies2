package movies

import (
	"context"
	"micromovies2/movies/pb"
	"github.com/golang/protobuf/ptypes"
)

//Encode and Decode GetMovies Request and response
func EncodeGRPCGetMoviesRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, nil
}

func DecodeGRPCGetMoviesRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, nil
}

// Encode and Decode GetMovies Response
func EncodeGRPCGetMoviesResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(getMoviesResponse)
	var movies []*pb.Movie
	for _, movie := range resp.Movies {
		createdOn, err := ptypes.TimestampProto(movie.CreatedOn)
		if err != nil {
			//todo bring logger
			return nil, err
		}
		updatedOn, err := ptypes.TimestampProto(movie.UpdatedOn)
		if err != nil {
			//todo bring logger
			return nil, err
		}
		//pb director type conversion
		var director []*pb.Director
		for _, d := range movie.Director {
			director = append(director, &pb.Director{Director: d})
		}
		m := &pb.Movie{
			Id:        movie.Id,
			Title:     movie.Title,
			Director:  director,
			Year:      movie.Year,
			Createdby: movie.CreatedBy,
			Createdon: createdOn,
			Updatedby: movie.UpdatedBy,
			Updatedon: updatedOn,
		}
		movies = append(movies, m)
	}
	return &pb.GetMoviesResponse{
		Movies: movies,
		Err:    resp.Err,
	}, nil
}

func DecodeGRPCGetMoviesResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.GetMoviesResponse)
	var movies []Movie
	for _, movie := range resp.Movies {
		createdOn, err := ptypes.Timestamp(movie.Createdon)
		if err != nil {
			//todo log error
			return nil, err
		}
		updatedOn, err := ptypes.Timestamp(movie.Updatedon)
		if err != nil {
			//todo log error
			return nil, err
		}
		var director []string
		for _, d := range movie.Director {
			director = append(director, d.Director)
		}
		m := Movie{
			Id:        movie.Id,
			Title:     movie.Title,
			Director:  director,
			Year:      movie.Year,
			CreatedBy: movie.Createdby,
			CreatedOn: createdOn,
			UpdatedBy: movie.Updatedby,
			UpdatedOn: updatedOn,
		}
		movies = append(movies, m)
	}
	return getMoviesResponse{
		Movies: movies,
		Err:    resp.Err,
	}, nil
}

//encode GetMoviesByIdRequest
func EncodeGRPCGetMovieByIdRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(getMovieByIdRequest)
	return &pb.GetMovieByIdRequest{
		Id: req.Id,
	}, nil
}

//decode GetMovieByIdRequest
func DecodeGRPCGetMovieByIdRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetMovieByIdRequest)
	return getMovieByIdRequest{
		Id: req.Id,
	}, nil
}

// encode GetMovieById Response
func EncodeGRPCGetMovieByIdResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(getMovieByIdResponse)
	createdOn, err := ptypes.TimestampProto(resp.Movie.CreatedOn)
	if err != nil {
		//todo bring logger
		return nil, err
	}
	updatedOn, err := ptypes.TimestampProto(resp.Movie.UpdatedOn)
	if err != nil {
		//todo bring logger
		return nil, err
	}
	//pb director type conversion
	var director []*pb.Director
	for _, d := range resp.Movie.Director {
		director = append(director, &pb.Director{Director: d})
	}
	m := &pb.Movie{
		Id:        resp.Movie.Id,
		Title:     resp.Movie.Title,
		Director:  director,
		Year:      resp.Movie.Year,
		Createdby: resp.Movie.CreatedBy,
		Createdon: createdOn,
		Updatedby: resp.Movie.UpdatedBy,
		Updatedon: updatedOn,
	}

	return &pb.GetMovieByIdResponse{
		Movie: m,
		Err:   resp.Err,
	}, nil
}

// decode GetMovieById Response
func DecodeGRPCGetMovieByIdResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.GetMovieByIdResponse)
	createdOn, err := ptypes.Timestamp(resp.Movie.Createdon)
	if err != nil {
		//todo log error
		return nil, err
	}
	updatedOn, err := ptypes.Timestamp(resp.Movie.Updatedon)
	if err != nil {
		//todo log error
		return nil, err
	}
	var director []string
	for _, d := range resp.Movie.Director {
		director = append(director, d.Director)
	}
	m := Movie{
		Id:        resp.Movie.Id,
		Title:     resp.Movie.Title,
		Director:  director,
		Year:      resp.Movie.Year,
		CreatedBy: resp.Movie.Createdby,
		CreatedOn: createdOn,
		UpdatedBy: resp.Movie.Updatedby,
		UpdatedOn: updatedOn,
	}

	return getMovieByIdResponse{
		Movie: m,
		Err:   resp.Err,
	}, nil
}

//encode NewMovieRequest
func EncodeGRPCNewMovieRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(newMovieRequest)
	//pb director type conversion
	var director []*pb.Director
	for _, d := range req.Director {
		director = append(director, &pb.Director{Director: d})
	}

	return &pb.NewMovieRequest{
		Title:     req.Title,
		Director:  director,
		Year:      req.Year,
		Createdby: req.Createdby,
	}, nil
}

//decode NewMovieRequest
func DecodeGRPCNewMovieRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.NewMovieRequest)
	var director []string
	for _, d := range req.Director {
		director = append(director, d.Director)
	}
	return newMovieRequest{
		Title:     req.Title,
		Director:  director,
		Year:      req.Year,
		Createdby: req.Createdby,
	}, nil
}

// Encode and Decode NewMovieResponse
func EncodeGRPCNewMovieResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(newMovieResponse)
	return &pb.NewMovieResponse{
		Id:  resp.Id,
		Err: resp.Err,
	}, nil
}

func DecodeGRPCNewMovieResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.NewMovieResponse)
	return newMovieResponse{
		Id:  resp.Id,
		Err: resp.Err,
	}, nil
}

//encode deleteMovieRequest
func EncodeGRPCDeleteMovieRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(deleteMovieRequest)
	return &pb.DeleteMovieRequest{
		Id: req.Id,
	}, nil
}

//decode DeleteMovieRequest
func DecodeGRPCDeleteMovieRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.DeleteMovieRequest)
	return deleteMovieRequest{
		Id: req.Id,
	}, nil
}

// Encode and Decode DeleteMovieResponse
func EncodeGRPCDeleteMovieResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(deleteMovieResponse)
	return &pb.DeleteMovieResponse{
		Err: resp.Err,
	}, nil
}

func DecodeGRPCDeleteMovieResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.DeleteMovieResponse)
	return deleteMovieResponse{
		Err: resp.Err,
	}, nil
}

//encode updateMovieRequest
func EncodeGRPCUpdateMovieRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(updateMovieRequest)
	var director []*pb.Director
	for _, d := range req.Director {
		director = append(director, &pb.Director{Director: d})
	}
	return &pb.UpdateMovieRequest{
		Id:        req.Id,
		Title:     req.Title,
		Director:  director,
		Year:      req.Year,
		Createdby: req.Createdby,
	}, nil
}

//decode UpdateMovieRequest
func DecodeGRPCUpdateMovieRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateMovieRequest)
	var director []string
	for _, d := range req.Director {
		director = append(director, d.Director)
	}
	return updateMovieRequest{
		Id:        req.Id,
		Title:     req.Title,
		Director:  director,
		Year:      req.Year,
		Createdby: req.Createdby,
	}, nil
}

// Encode and Decode UpdateMovieResponse
func EncodeGRPCUpdateMovieResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(updateMovieResponse)
	return &pb.UpdateMovieResponse{
		Err: resp.Err,
	}, nil
}

func DecodeGRPCUpdateMovieResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UpdateMovieResponse)
	return updateMovieResponse{
		Err: resp.Err,
	}, nil
}
