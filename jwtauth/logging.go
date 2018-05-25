package jwtauth

import (
	"github.com/rs/zerolog"
	"time"
	"context"
)

//struct passing the logger
type LoggingMiddleware struct {
	Logger zerolog.Logger
	Next Service
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) GenerateToken(ctx context.Context, email string, role string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "GenerateToken").Err(err).Dur("took", time.Since(begin)).
			Str("email", email).Str("role",role).Msg("")
	}(time.Now())
	output, err = mw.Next.GenerateToken(ctx, email, role)
	return
}
