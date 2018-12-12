package main

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"mime"
	"os"
	"path/filepath"
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

func followRedirect(u string) (string, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(u)

	if err := fasthttp.Do(req, res); err != nil {
		return "", err
	}

	if sc := res.StatusCode(); sc != 302 {
		return "", fmt.Errorf("failed to get redirected to mp3: HTTP status %d", sc)
	}

	return string(res.Header.Peek("Location")), nil
}

func download(u string) (int64, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(u)

	if err := fasthttp.Do(req, res); err != nil {
		return 0, err
	}

	if sc := res.StatusCode(); sc != 200 {
		return 0, fmt.Errorf("HTTP status %d", sc)
	}

	cd := string(res.Header.Peek("Content-Disposition"))
	if cd == "" {
		return 0, fmt.Errorf("missing Content-Disposition header")
	}

	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return 0, err
	}

	fileName := params["filename"]
	if fileName == "" {
		return 0, fmt.Errorf("missing file name in Content-Disposition header")
	}

	fileName = filepath.Join(*dir, fileName)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return 0, err
	}

	return res.WriteTo(f)
}
