package users

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type LoggingMiddleware struct {
	Logger zap.Logger
	Next   Service
}

func (mw LoggingMiddleware) NewUser(ctx context.Context, user User) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "NewUser"), zap.String("name", user.Name), zap.String("lastname", user.LastName),
			zap.String("email", user.Email), zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.NewUser(ctx, user)
	return
}

func (mw LoggingMiddleware) GetUserByEmail(ctx context.Context, email string) (output User, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "GetUserByEmail"), zap.String("email", email),
			zap.Error(err), zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.GetUserByEmail(ctx, email)
	return
}

func (mw LoggingMiddleware) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (output bool, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "ChangePassword"), zap.String("email", email), zap.Error(err),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.ChangePassword(ctx, email, currentPassword, newPassword)
	return
}

func (mw LoggingMiddleware) Login(ctx context.Context, email string, Password string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info("", zap.String("method", "Login"), zap.String("email", email), zap.Error(err),
			zap.Duration("took", time.Since(begin)))
	}(time.Now())
	output, err = mw.Next.Login(ctx, email, Password)
	return
}
