package worker

import (
	"context"
	"fmt"
	"time"
)

type IntervalWorker struct {
	Interval time.Duration
	period time.Duration
	job Job
}

func NewIntervalWorker(interval time.Duration) *IntervalWorker {
	return &IntervalWorker{
		Interval: interval,
		period: interval,
	}
}

func (iw *IntervalWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Cancelled interval worker. Error detail :%v\n", ctx.Err())
			return
		case <-time.After(iw.period):
		}

		started := time.Now()
		iw.execute(ctx)
		finished := time.Now()

		duration := finished.Sub(started)
		iw.period = iw.Interval - duration
	}
}

func (iw *IntervalWorker) execute(ctx context.Context) {
	iw.job.execute(ctx)
}