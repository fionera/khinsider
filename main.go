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
var numGames int64
var queuedJobs sync.WaitGroup
var crawlerGroup sync.WaitGroup
var availableLetters []*Letter
var exitRequested int32

var tracks chan Job
var files chan Job
var games chan Job
var letters chan Job

type Job interface {
	Crawl(c context.Context) error
}

func init() {
	l := []*Letter{{Letter: '#'}}
	for c := 'A'; c <= 'Z'; c++ {
		l = append(l, &Letter{Letter: c})
	}
	availableLetters = l
}

func main() {
	if err := parseArgs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting Kingdom Hearts Insider Scraper")
	logrus.Info("  https://github.com/fionera/khinsider/")

	c, cancel := context.WithCancel(context.Background())

	go listenCtrlC(cancel)
	go stats()

	letters = make(chan Job, *concurrency)
	crawlerGroup.Add(int(*concurrency))
	for i := 0; i < int(*concurrency); i++ {
		go crawler(c, letters)
	}

	games = make(chan Job, *concurrency)
	crawlerGroup.Add(int(*concurrency))
	for i := 0; i < int(*concurrency); i++ {
		go crawler(c, games)
	}

	tracks = make(chan Job, *concurrency)
	crawlerGroup.Add(int(*concurrency))
	for i := 0; i < int(*concurrency); i++ {
		go crawler(c, tracks)
	}

	files = make(chan Job, *concurrency)
	crawlerGroup.Add(int(*concurrency))
	for i := 0; i < int(*concurrency); i++ {
		go crawler(c, files)
	}

	// Start letter grabber
	for _, letter := range availableLetters {
		queuedJobs.Add(1)
		letters <- letter
	}
	close(letters)

	// Shutdown
	queuedJobs.Wait()
	close(games)
	close(tracks)
	close(files)
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
			"games":       numGames,
			"tracks":      numDownloaded,
			"total_bytes": totalBytes,
			"avg_rate":    fmt.Sprintf("%.0f", float64(total)/dur),
		}).Info("Stats")
	}
}
