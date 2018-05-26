package vault

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"time"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           Service
}

func (mw InstrumentingMiddleware) Hash(ctx context.Context, password string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "hash", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Next.Hash(ctx, password)
	return
}

func (mw InstrumentingMiddleware) Validate(ctx context.Context, password string, hash string) (bool, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "count", "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return mw.Next.Validate(ctx, password, hash)

}
