package vault

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
func (mw LoggingMiddleware) Hash(ctx context.Context, password string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str(
			"method", "hash").Str("password", password).Str("output", output).Err(err).Dur("took",
			time.Since(begin)).Msg("")

	}(time.Now())
	output, err = mw.Next.Hash(ctx, password)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) Validate(ctx context.Context, password string, hash string) (output bool, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "validate").Str("password", password).Str("hash", hash).Bool("output",
			output).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.Validate(ctx, password, hash)
	return
}
