package retry

import (
	"time"
)

type Retry struct {
	maxRetry         int
	retryDurationVec []time.Duration
}

var DefaultHttpRetryVec = []time.Duration{
	time.Millisecond * 200,
	time.Millisecond * 200,
	time.Millisecond * 200,
}

func New(maxRetry int, retryDurationVec ...time.Duration) *Retry {
	if len(retryDurationVec) == 0 {
		retryDurationVec = DefaultHttpRetryVec
	}
	return &Retry{
		maxRetry:         maxRetry,
		retryDurationVec: retryDurationVec,
	}
}

func (r *Retry) Run(fn func() bool) bool {
	for i := 0; i <= r.maxRetry; i++ {
		if fn() {
			return true
		}

		if i == r.maxRetry {
			return false
		}

		sleepTime := time.Second
		durationLen := len(r.retryDurationVec)
		if r.retryDurationVec != nil && durationLen > 0 {
			if i < durationLen {
				sleepTime = r.retryDurationVec[i]
			} else {
				sleepTime = r.retryDurationVec[durationLen-1]
			}
		}
		select {
		case <-time.After(sleepTime):
		}
	}
	return false
}
