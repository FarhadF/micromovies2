package jwtauth

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
func (mw LoggingMiddleware) GenerateToken(ctx context.Context, email string, role string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "GenerateToken"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.Error(err), zap.String("email", email), zap.String(
			"role", role), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.GenerateToken(ctx, email, role)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) ParseToken(ctx context.Context, token string) (output Claims, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "ParseToken"), zap.String("correlationid",
			ctx.Value("correlationid").(string)), zap.Error(err), zap.String("token", token),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.ParseToken(ctx, token)
	return
}
