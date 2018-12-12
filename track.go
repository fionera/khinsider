package main

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"log"
)

type Track struct {
	Game     Game
	Title    string
	Download string
}

func (t *Track) Crawl(c context.Context) error {

	logrus.Println("https://downloads.khinsider.com" + t.Download)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI("https://downloads.khinsider.com" + t.Download)

	if err := fasthttp.Do(req, res); err != nil {
		return err
	}

	logrus.Infof("Visited Song %s", t.Title)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#EchoTopic > audio").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		src, exist := s.Attr("src")

		if exist {
			size, err := download(src)
			if err != nil {
				log.Fatal(err)
			}

			totalBytes += size
		}
	})

	return nil
}
