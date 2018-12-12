package main

import (
	"context"
	"fmt"
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

	logrus.Infof("Visited letter %s", string(l.Letter))

	//tracks := readGames(bytes.NewReader(res.Body()))
	//for _, track := range tracks {
	//	select {
	//	case <-c.Done():
	//		return nil
	//	case jobs <- track:
	//		continue
	//	}
	//
	//}

	return nil
}
