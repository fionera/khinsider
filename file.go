package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
	"sync/atomic"
)

type File struct {
	Track Track
	URL   []byte
}

func (f *File) Crawl(c context.Context) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(string(f.URL))

	if err := fasthttp.Do(req, res); err != nil {

		return nil
	}

	if sc := res.StatusCode(); sc != 200 {
		return fmt.Errorf("HTTP status %d", sc)
	}

	folder := filepath.Join(*dir, string(f.Track.Game.Name))

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return nil
	}

	fileName := filepath.Join(folder, filepath.Base(string(f.URL)))

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		logrus.Infof("Skipping File " + fileName)
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil
	}

	size, err := res.WriteTo(file)
	if err != nil {
		return nil
	}

	atomic.AddInt64(&totalBytes, size)
	atomic.AddInt64(&numDownloaded, 1)

	return file.Close()
}
