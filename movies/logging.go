package movies

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

//struct passing the logger
type LoggingMiddleware struct {
	Logger zerolog.Logger
	Next   Service
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) GetMovies(ctx context.Context) (output []Movie, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "GetMovies").Str("correlationid",
			ctx.Value("correlationid").(string)).Err(err).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.GetMovies(ctx)
	return
}

func (mw LoggingMiddleware) GetMovieById(ctx context.Context, id string) (output Movie, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "GetMovieById").Str("correlationid",
			ctx.Value("correlationid").(string)).Str("id", id).
			Err(err).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.GetMovieById(ctx, id)
	return
}

func (mw LoggingMiddleware) NewMovie(ctx context.Context, title string, director []string, year string, userid string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "NewMovie").Str("correlationid",
			ctx.Value("correlationid").(string)).Str("title", title).Strs("director", director).
			Err(err).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.NewMovie(ctx, title, director, year, userid)
	return
}

func (mw LoggingMiddleware) DeleteMovie(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "DeleteMovie").Str("correlationid",
			ctx.Value("correlationid").(string)).Str("id", id).Err(err).Dur("took",
			time.Since(begin)).Msg("")
	}(time.Now())
	err = mw.Next.DeleteMovie(ctx, id)
	return
}

func (mw LoggingMiddleware) UpdateMovie(ctx context.Context, id string, title string, director []string, year string, userid string) (err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "UpdateMovie").Str("correlationid",
			ctx.Value("correlationid").(string)).Str("id", id).Err(err).Dur("took",
			time.Since(begin)).Msg("")
	}(time.Now())
	err = mw.Next.UpdateMovie(ctx, id, title, director, year, userid)
	return
}
