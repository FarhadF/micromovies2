package apigateway

import (
	"context"
	"github.com/farhadf/micromovies2/movies"
	"go.uber.org/zap"
	"time"
)

//struct passing the logger
type LoggingMiddleware struct {
	Logger zap.Logger
	Next   Service
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) Login(ctx context.Context, email string, password string) (token string, refreshToken string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "Login"),
			zap.String("correlationid", ctx.Value("correlationid").(string)), zap.Error(err),
			zap.String("email", email), zap.String("password", password), zap.Duration("took",
				time.Since(begin)))
	}(time.Now())
	token, refreshToken, err = mw.Next.Login(ctx, email, password)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) Register(ctx context.Context, email string, password string, firstname string, lastname string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "Register"),
			zap.String("correlationid", ctx.Value("correlationid").(string)), zap.Error(err),
			zap.String("email", email), zap.String("password", password), zap.String("firstname", firstname),
			zap.String("lastname", lastname), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.Register(ctx, email, password, firstname, lastname)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (output bool, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "ChangePassword"),
			zap.String("correlationid", ctx.Value("correlationid").(string)),
			zap.String("email", email), zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.ChangePassword(ctx, email, currentPassword, newPassword)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) GetMovieById(ctx context.Context, id string) (output movies.Movie, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "GetMovieById"),
			zap.String("correlationid", ctx.Value("correlationid").(string)),
			zap.String("id", id), zap.Error(err),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.GetMovieById(ctx, id)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) NewMovie(ctx context.Context, title string, director []string, year string, userId string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "NewMovie"),
			zap.String("correlationid", ctx.Value("correlationid").(string)),
			zap.String("title", title),
			zap.String("userId", userId), zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.NewMovie(ctx, title, director, year, userId)
	return
}
