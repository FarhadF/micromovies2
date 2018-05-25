package jwtauth

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
func (mw InstrumentingMiddleware) GenerateToken(ctx context.Context, email string, role string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GenerateToken", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.GenerateToken(ctx, email, role)
	return
}
