package main

import (
	"context"
	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

func crawler(c context.Context) {
	defer crawlerGroup.Done()

	for {
		if atomic.LoadInt32(&exitRequested) != 0 {
			break
		}

		select {
		case job := <-jobs:
			err := backoff.Retry(func() error {
				err := job.Crawl(c)
				if err != nil {
					logrus.WithError(err).
						Errorf("Failed crawling")
				}
				return err
			}, backoff.NewExponentialBackOff())

			if err != nil {
				logrus.Fatal(err)
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
