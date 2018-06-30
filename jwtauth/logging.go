package jwtauth

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/satori/go.uuid"
	"time"
)

//struct passing the logger
type LoggingMiddleware struct {
	Logger zerolog.Logger
	Next   Service
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) GenerateToken(ctx context.Context, email string, role string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "GenerateToken").Str("correlationid",
			ctx.Value("correlationid").(uuid.UUID).String()).Err(err).Dur("took", time.Since(begin)).
			Str("email", email).Str("role", role).Msg("")
	}(time.Now())
	output, err = mw.Next.GenerateToken(ctx, email, role)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) ParseToken(ctx context.Context, token string) (output Claims, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "ParseToken").Str("correlationid",
			ctx.Value("correlationid").(uuid.UUID).String()).Err(err).Dur("took", time.Since(begin)).
			Str("token", token).Msg("")
	}(time.Now())
	output, err = mw.Next.ParseToken(ctx, token)
	return
}
