package main

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"log"
	"sync/atomic"
)

type Track struct {
	Game     Game
	Title    []byte
	Download []byte
}

func (t *Track) Crawl(c context.Context) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI("https://downloads.khinsider.com" + string(t.Download))

	if err := fasthttp.Do(req, res); err != nil {
		return err
	}

	logrus.Infof("Visited Song | %s - %s", t.Game.Name, t.Title)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#EchoTopic > audio").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		src, exist := s.Attr("src")

		if exist {
			size, err := download(src, *t)
			if err != nil {
				log.Fatal(err)
			}

			atomic.AddInt64(&totalBytes, size)
			if size != 0 {
				atomic.AddInt64(&numDownloaded, 1)
			}
		}
	})

	return nil
}
