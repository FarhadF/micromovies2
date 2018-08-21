package users

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"time"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	//CountResult    metrics.Histogram
	Next Service
}

//instrumentation per method
func (mw InstrumentingMiddleware) NewUser(ctx context.Context, user User) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewUser", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.NewUser(ctx, user)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) GetUserByEmail(ctx context.Context, email string) (output User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserByEmail", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.GetUserByEmail(ctx, email)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) ChangePassword(ctx context.Context, email string, currentPassword string, newPassword string) (output bool, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ChangePassword", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.ChangePassword(ctx, email, currentPassword, newPassword)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) Login(ctx context.Context, email string, Password string) (token string, refreshToken string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	token, refreshToken, err = mw.Next.Login(ctx, email, Password)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) Refresh(ctx context.Context, token string, refreshToken string) (newToken string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Refresh", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	token, err = mw.Next.Refresh(ctx, token, refreshToken)
	return
}
