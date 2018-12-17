package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Letter struct {
	Letter []byte
	Games  []Game
}

func (l *Letter) Crawl(c context.Context) error {
	u := fmt.Sprintf("https://downloads.khinsider.com/game-soundtracks/browse/%s", string(l.Letter))

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(u)

	if err := fasthttp.Do(req, res); err != nil {
		return err
	}

	logrus.Debugf("Visited letter | %s", string(l.Letter))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		logrus.Error(err)
		return err
	}

	doc.Find("#EchoTopic > p:nth-child(5) > a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, exist := s.Attr("href")
		name := s.Text()

		if exist {
			jobs <- &Game{
				URL:  []byte(href),
				Name: []byte(name),
			}
		}
	})

	return err
}
