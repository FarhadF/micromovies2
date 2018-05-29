package apigateway

import (
	"context"
	"go.uber.org/zap"
	"time"
)

//struct passing the logger
type LoggingMiddleware struct {
	Logger zap.Logger
	Next   Service
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) Login(ctx context.Context, email string, password string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "Login"), zap.Error(err),
			zap.Duration("took", time.Since(begin)), zap.String("email", email),
			zap.String("password", password))
	}(time.Now())
	output, err = mw.Next.Login(ctx, email, password)
	return
}

//each method will have its own logger for app logs
func (mw LoggingMiddleware) Register(ctx context.Context, email string, password string, firstname string, lastname string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "Register"), zap.Error(err),
			zap.Duration("took", time.Since(begin)), zap.String("email", email),
			zap.String("password", password), zap.String("firstname", firstname),
			zap.String("lastname", lastname))
	}(time.Now())
	output, err = mw.Next.Register(ctx, email, password, firstname, lastname)
	return
}
