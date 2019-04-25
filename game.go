package main

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"sync/atomic"
)

type Game struct {
	Letter Letter
	URL    []byte
	Name   []byte
}

func (g *Game) Crawl(c context.Context) error {
	defer queuedJobs.Done()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI("https://downloads.khinsider.com" + string(g.URL))

	if err := fasthttp.Do(req, res); err != nil {
		return err
	}

	logrus.Debugf("Visited Game | %s", g.Name)
	atomic.AddInt64(&numGames, 1)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		logrus.Error(err)
		return err
	}

	doc.Find("#songlist > tbody > tr > td:nth-child(3) > a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, exist := s.Attr("href")
		name := s.Text()

		if exist {
			queuedJobs.Add(1)

			tracks <- &Track{
				Game:     *g,
				Download: []byte(href),
				Title:    []byte(name),
			}
		}
	})

	return err
}
