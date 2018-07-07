package movies

import (
	"context"

	"errors"
	"github.com/jackc/pgx"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"reflect"
	"time"
)

type Service interface {
	GetMovies(ctx context.Context) ([]Movie, error)
	GetMovieById(ctx context.Context, id string) (Movie, error)
	NewMovie(ctx context.Context, title string, director []string, year string, createdBy string) (string, error)
	DeleteMovie(ctx context.Context, id string) error
	UpdateMovie(ctx context.Context, id string, title string, director []string, year string, updatedBy string) error
}

//implementation with database and logger
type moviesService struct {
	db     *pgx.ConnPool
	logger *zap.Logger
}

//constructor - we can later add initialization if needed
func NewService(db *pgx.ConnPool, logger *zap.Logger) Service {
	return moviesService{
		db,
		logger,
	}
}

//implementation
func (m moviesService) GetMovies(ctx context.Context) ([]Movie, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("GetMovies", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	rows, err := m.db.Query("select * from movies")
	if err != nil {
		return nil, err
	}
	movies := make([]Movie, 0)
	for rows.Next() {
		movie := new(Movie)
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.CreatedBy, &movie.CreatedOn, &movie.UpdatedBy, &movie.UpdatedOn)
		if err != nil {
			return nil, err
		}

		r, err := m.db.Query("select director from movie_directors where movie_id=$1", movie.Id)
		if err != nil {
			return nil, err
		}
		var director []string
		var d string
		for r.Next() {
			err = r.Scan(&d)
			if err != nil {
				return nil, err
			}
			director = append(director, d)
		}
		movie.Director = director
		movies = append(movies, *movie)
	}
	return movies, nil
}

//implementation
func (m moviesService) GetMovieById(ctx context.Context, id string) (Movie, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("GetMovieById", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		span.SetTag("id", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	var movie Movie
	row := m.db.QueryRow("select * from movies where id = $1", id)
	err := row.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.CreatedBy, &movie.CreatedOn, &movie.UpdatedBy, &movie.UpdatedOn)
	if err != nil {
		return movie, err
	}
	r, err := m.db.Query("select director from movie_directors where movie_id=$1", movie.Id)
	var director []string
	for r.Next() {
		var d string
		err = r.Scan(&d)
		if err != nil {
			return movie, err
		}
		director = append(director, d)
	}
	movie.Director = director
	return movie, nil
}

//implementation
//todo: check if the role is user and if so discard createdBy if available and extract email from ctx
func (m moviesService) NewMovie(ctx context.Context, title string, director []string, year string, createdBy string) (string, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("NewMovie", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	rows, err := m.db.Query("select * from movies where title=$1", title)
	defer rows.Close()
	if err != nil {
		//todo: add logging
		return "", err
	}
	if !rows.Next() {
		var id string
		err := m.db.QueryRow("insert into movies (title, year, createdBy) values($1,$2,$3) returning id", title, year, createdBy).Scan(&id)
		//res, err := stmt.Exec(movie.Title,movie.Director, movie.Year, movie.Userid)
		//id, err := res.LastInsertId()
		if err != nil {
			return "", err
		}
		for _, d := range director {
			_, err = m.db.Exec("insert into movie_directors (movie_id, director, createdby) values($1,$2,$3)", id, d, createdBy)
			if err != nil {
				//rollback
				err1 := err
				_, err := m.db.Exec("delete from movies where id=$1", id)
				if err != nil {
					return "", err
				}
				return "", err1
			}
		}
		//return strconv.FormatInt(id, 10), nil
		return id, nil
	} else {

		return "", errors.New("movie already exists")
	}
}

//implementation
func (m moviesService) DeleteMovie(ctx context.Context, id string) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("DeleteMovie", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		span.SetTag("id", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	rows, err := m.db.Query("select * from movies where id=$1", id)
	defer rows.Close()
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("movie does not exist")
	}
	_, err = m.db.Exec("delete from movies where id = $1", id)
	if err != nil {
		return err
	}
	_, err = m.db.Exec("delete from movie_directors where movie_id = $1", id)
	return nil
}

//implementation
func (m moviesService) UpdateMovie(ctx context.Context, id string, title string, director []string, year string, updatedBy string) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("UpdateMovie", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		span.SetTag("id", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	rows, err := m.db.Query("select * from movies where id=$1", id)
	defer rows.Close()
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("movie does not exist")
	}
	updatedOn := time.Now().UTC()
	_, err = m.db.Exec("update movies set title = $1, year = $2, updatedon = $3, updatedby = $4 where id = $5", title, year,
		updatedOn.Format("2006-01-02 15:04:05.999999"), updatedBy, id)
	if err != nil {
		return err
	}
	r, err := m.db.Query("select director from movie_directors where movie_id=$1", id)
	defer r.Close()
	if err != nil {
		//todo rollback
		return err
	}
	var dir []string
	for r.Next() {
		var d string
		err = r.Scan(&d)
		if err != nil {
			return err
		}
		dir = append(dir, d)
	}
	if reflect.DeepEqual(dir, director) {
		return nil
	}
	_, err = m.db.Exec("delete from movie_directors where movie_id = $1", id)
	if err != nil {
		//todo: rollback
		return err
	}
	for _, d := range director {
		_, err = m.db.Exec("insert into movie_directors (movie_id, director, createdby, createdOn) values($1,$2,$3,$4)",
			id, d, updatedBy, updatedOn)
		if err != nil {
			//todo:rollback
			return err
		}
	}
	return nil
}
