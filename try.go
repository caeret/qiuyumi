package qiuyumi

import (
	"math/rand"
	"time"
)

type stop struct {
	error
}

func Try(attempts int, sleep time.Duration, fn func() error) (err error) {
	if err = fn(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return Try(attempts, 2*sleep, fn)
		}
	}
	return
}
