package apigateway

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
func (mw InstrumentingMiddleware) Login(ctx context.Context, email string, password string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.Login(ctx, email, password)
	return
}
