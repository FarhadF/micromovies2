package users

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type LoggingMiddleware struct {
	Logger zerolog.Logger
	Next   Service
}

func (mw LoggingMiddleware) NewUser(ctx context.Context, user User) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "NewUser").Str("name", user.Name).Str("lastname", user.LastName).
			Str("email", user.Email).Err(err).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.NewUser(ctx, user)
	return
}

func (mw LoggingMiddleware) GetUserByEmail(ctx context.Context, email string) (output User, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "GetUserByEmail").
			Str("email", email).Err(err).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.GetUserByEmail(ctx, email)
	return
}

func (mw LoggingMiddleware) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (output bool, err error) {
	defer func(begin time.Time) {
		mw.Logger.Info().Str("method", "ChangePassword").
			Str("email", email).Err(err).Dur("took", time.Since(begin)).Msg("")
	}(time.Now())
	output, err = mw.Next.ChangePassword(ctx, email, currentPassword, newPassword)
	return
}
