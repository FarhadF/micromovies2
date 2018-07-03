package movies

import (
	"context"
	"go.uber.org/zap"
	"time"
)

//struct passing the logger
type LoggingMiddleware struct {
	Logger *zap.Logger
	Next   Service
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) GetMovies(ctx context.Context) (output []Movie, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "GetMovies"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.GetMovies(ctx)
	return
}

func (mw LoggingMiddleware) GetMovieById(ctx context.Context, id string) (output Movie, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "GetMovieById"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.String("id", id), zap.Error(err),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.GetMovieById(ctx, id)
	return
}

func (mw LoggingMiddleware) NewMovie(ctx context.Context, title string, director []string, year string, userid string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "NewMovie"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.String("title", title), zap.Strings("director", director),
			zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.NewMovie(ctx, title, director, year, userid)
	return
}

func (mw LoggingMiddleware) DeleteMovie(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "DeleteMovie"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.String("id", id), zap.Error(err),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	err = mw.Next.DeleteMovie(ctx, id)
	return
}

func (mw LoggingMiddleware) UpdateMovie(ctx context.Context, id string, title string, director []string, year string, userid string) (err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "UpdateMovie"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.String("id", id), zap.Error(err),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	err = mw.Next.UpdateMovie(ctx, id, title, director, year, userid)
	return
}
