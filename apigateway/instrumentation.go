package apigateway

import (
	"context"
	"fmt"
	"github.com/farhadf/micromovies2/movies"
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
func (mw InstrumentingMiddleware) Login(ctx context.Context, email string, password string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.Login(ctx, email, password)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) Register(ctx context.Context, email string, password string, firstname string, lastname string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.Register(ctx, email, password, firstname, lastname)
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
func (mw InstrumentingMiddleware) GetMovieById(ctx context.Context, id string) (output movies.Movie, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetMovieById", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.GetMovieById(ctx, id)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) NewMovie(ctx context.Context, title string, director []string, year string, userId string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewMovie", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.NewMovie(ctx, title, director, year, userId)
	return
}
