package users

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
func (mw InstrumentingMiddleware) NewUser(ctx context.Context, user User) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewUser", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.Next.NewUser(ctx, user)
	return
}
