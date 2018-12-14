package main

import (
	"context"
	"fmt"
	"sync/atomic"
)

func crawler(c context.Context) {
	defer crawlerGroup.Done()

	for {
		if atomic.LoadInt32(&exitRequested) != 0 {
			break
		}

		select {
		case job := <-jobs:
			if err := job.Crawl(c); err != nil {
				fmt.Println(err)
			}
		default:
		}
	}
}
