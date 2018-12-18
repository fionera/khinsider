package main

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"path/filepath"
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

	logrus.Debugf("Visited Song | %s - %s", t.Game.Name, t.Title)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		logrus.Error(err)
		return err
	}

	doc.Find("#EchoTopic > p > a[href*='/ost/']").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		src, exist := s.Attr("href")

		logrus.Debugf("Found File | %s - %s - %s", t.Game.Name, t.Title, filepath.Ext(src))

		if exist {
			jobs <- &File{
				Track: *t,
				URL:   []byte(src),
			}
		}
	})

	return nil
}
