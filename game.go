package main

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"log"
)

type Game struct {
	URL  string
	Name string
}

func (g *Game) Crawl(c context.Context) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI("https://downloads.khinsider.com" + g.URL)

	if err := fasthttp.Do(req, res); err != nil {
		return err
	}

	logrus.Infof("Visited Game %s", g.Name)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#songlist > tbody > tr > td:nth-child(3) > a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, exist := s.Attr("href")
		name := s.Text()

		if exist {
			jobs <- &Track{
				Game:     *g,
				Download: href,
				Title:    name,
			}
		}
	})

	return err
}
