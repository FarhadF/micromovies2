package vault

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
func (mw LoggingMiddleware) Hash(ctx context.Context, password string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "hash"), zap.String("password", password),
			zap.String("correlationid", ctx.Value("correlationid").(string)), zap.String("output", output),
			zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.Hash(ctx, password)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) Validate(ctx context.Context, password string, hash string) (output bool, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "validate"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.String("password", password),
			zap.String("hash", hash), zap.Bool("output", output), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.Validate(ctx, password, hash)
	return
}
