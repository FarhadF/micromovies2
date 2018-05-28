package apigateway

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
func (mw LoggingMiddleware) Login(ctx context.Context, email string, password string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "Login").Err(err).Dur("took", time.Since(begin)).
			Str("email", email).Str("password", password).Msg("")
	}(time.Now())
	output, err = mw.Next.Login(ctx, email, password)
	return
}
