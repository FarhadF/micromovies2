package movies

import (
	"github.com/go-kit/kit/metrics"
	"time"
	"fmt"
	"context"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	//CountResult    metrics.Histogram
	Next           Service
}

//instrumentation per method
func (mw InstrumentingMiddleware) GetMovies(ctx context.Context) (output []Movie, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetMovies", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.GetMovies(ctx)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) GetMovieById(ctx context.Context, id string) (output Movie, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetMovieById", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.GetMovieById(ctx, id)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) NewMovie(ctx context.Context, title string, director []string, year string, userid string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewMovie", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.NewMovie(ctx, title, director, year, userid)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) UpdateMovie(ctx context.Context, id string, title string, director []string, year string, userid string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateMovie", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Next.UpdateMovie(ctx, id, title, director, year, userid)
	return
}

//instrumentation per method
func (mw InstrumentingMiddleware) DeleteMovie(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteMovie", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Next.DeleteMovie(ctx, id)
	return
}
