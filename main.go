package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var startTime = time.Now()
var totalBytes int64
var numDownloaded int64
var crawlerGroup sync.WaitGroup
var exitRequested int32
var jobs chan Job

type Job interface {
	Crawl(c context.Context) error
}

func main() {
	if err := parseArgs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logrus.Info("Starting Kingdom Hearts Insider Scraper")
	logrus.Info("  https://github.com/fionera/khinsider/")

	c, cancel := context.WithCancel(context.Background())

	jobs = make(chan Job)
	go listenCtrlC(cancel)
	go stats()

	// Start downloaders
	crawlerGroup.Add(int(*concurrency))
	for i := 0; i < int(*concurrency); i++ {
		go crawler(c)
	}

	// Start letter grabber
	for _, letter := range availableLetters {
		jobs <- &Letter{
			Letter: []byte(letter),
		}
	}

	// Shutdown
	close(jobs)
	crawlerGroup.Wait()

	total := atomic.LoadInt64(&totalBytes)
	dur := time.Since(startTime).Seconds()

	logrus.WithFields(logrus.Fields{
		"total_bytes": total,
		"dur":         dur,
		"avg_rate":    float64(total) / dur,
	}).Info("Stats")
}

func listenCtrlC(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	atomic.StoreInt32(&exitRequested, 1)
	cancel()
	fmt.Fprintln(os.Stderr, "\nWaiting for downloads to finish...")
	fmt.Fprintln(os.Stderr, "Press ^C again to exit instantly.")
	<-c
	fmt.Fprintln(os.Stderr, "\nKilled!")
	os.Exit(255)
}

func stats() {
	for range time.NewTicker(time.Second).C {
		total := atomic.LoadInt64(&totalBytes)
		dur := time.Since(startTime).Seconds()

		logrus.WithFields(logrus.Fields{
			"tracks":      numDownloaded,
			"total_bytes": totalBytes,
			"avg_rate":    fmt.Sprintf("%.0f", float64(total)/dur),
		}).Info("Stats")
	}
}
