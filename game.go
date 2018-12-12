package main

import (
	"context"
)

type Game struct {
	URL    string
	Name   string
	Tracks []Track
}

func (g *Game) Crawl(c context.Context) error {

	println(g.Name)

	//req := fasthttp.AcquireRequest()
	//defer fasthttp.ReleaseRequest(req)
	//res := fasthttp.AcquireResponse()
	//defer fasthttp.ReleaseResponse(res)
	//
	//req.SetRequestURI(g.URL)
	//
	//if err := fasthttp.Do(req, res); err != nil {
	//	return err
	//}
	//
	//logrus.Infof("Visited page %d", g.Name)
	//
	////tracks := readGames(bytes.NewReader(res.Body()), g)
	////for _, track := range tracks {
	////	select {
	////	case <-c.Done():
	////		return nil
	////	case jobs <- &track:
	////		continue
	////	}
	////}
	////
	////jobs <- &Track{
	////	Title:    "",
	////	Download: "",
	////	Game:     *g,
	////}

	return nil
}

func validText(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	if b[0] == ' ' || b[0] == '\n' || b[0] == ',' {
		return false
	}
	return true
}
