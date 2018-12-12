package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Game struct {
	URL    string
	Name   string
	Tracks []Track
}

func (g *Game) Crawl(c context.Context) error {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(g.URL)

	if err := fasthttp.Do(req, res); err != nil {
		return err
	}

	logrus.Infof("Visited page %d", g.Name)

	//tracks := readGames(bytes.NewReader(res.Body()), g)
	//for _, track := range tracks {
	//	select {
	//	case <-c.Done():
	//		return nil
	//	case jobs <- &track:
	//		continue
	//	}
	//}
	//
	//jobs <- &Track{
	//	Title:    "",
	//	Download: "",
	//	Game:     *g,
	//}

	return nil
}

//func readGames(r io.Reader, g Game) (tracks []Track) {
//	doc := html.NewTokenizer(r)
//
//	var track Track
//
//	// Parser state:
//	// -1: Start of File
//	//  0: Expecting non-text tag
//	//  1: Expecting text tag
//	//  2: Expecting download link
//	//  3: Expecting genre
//	var state = -1
//	var field *string
//
//	var tagName, k, v []byte
//	var hasAttr bool
//
//	for {
//		t := doc.Next()
//
//		if t == html.StartTagToken {
//			tagName, hasAttr = doc.TagName()
//			if !hasAttr { goto nextToken }
//
//			// Check for new item
//			if atom.Lookup(tagName) == atom.Div {
//				k, v, hasAttr = doc.TagAttr()
//				if atom.Lookup(k) == atom.Class && bytes.HasPrefix(v, []byte("play-item")) {
//					track = Track{}
//					state = 0
//				}
//				goto nextToken
//			}
//		}
//
//		switch {
//		// Expecting span meta element
//		case t == html.StartTagToken && state == 0:
//			if atom.Lookup(tagName) != atom.Span { goto nextToken }
//
//			var k, v []byte
//			k, v, hasAttr = doc.TagAttr()
//
//			if atom.Lookup(k) == atom.Class {
//				switch {
//				case bytes.Equal(v, []byte("ptxt-artist")):
//					field = &track.Artist
//					state = 1
//
//				case bytes.Equal(v, []byte("ptxt-track")):
//					field = &track.Title
//					state = 1
//
//				case bytes.Equal(v, []byte("ptxt-album")):
//					field = &track.Game
//					state = 1
//
//				case bytes.Equal(v, []byte("ptxt-genre")):
//					state = 3
//
//				case bytes.Equal(v, []byte("playicn")):
//					state = 2
//				}
//			}
//
//		// Expecting download link
//		case t == html.StartTagToken && state == 2:
//			if atom.Lookup(tagName) != atom.A { goto nextToken }
//
//			for {
//				k, v, hasAttr = doc.TagAttr()
//
//				if atom.Lookup(k) == atom.Href &&
//					bytes.HasPrefix(v, []byte("https://freemusicarchive.org/music/download/")){
//					track.Download = string(v)
//
//					// Commit token
//					tracks = append(tracks, track)
//
//					state = -1
//					goto nextToken
//				}
//
//				if !hasAttr { break }
//			}
//
//		// Expecting HTML text
//		case t == html.TextToken && state == 1:
//			text := doc.Text()
//			if !validText(text) { goto nextToken }
//			*field = string(text)
//			state = 0
//
//		// Expecting Genre
//		case t == html.TextToken && state == 3:
//			text := doc.Text()
//			if !validText(text) { goto nextToken }
//			if bytes.Equal(text, []byte(", ")) { goto nextToken }
//			track.Genres = append(track.Genres, string(text))
//
//		// Genre list closed
//		case t == html.EndTagToken && state == 3:
//			tagName, _ = doc.TagName()
//			if atom.Lookup(tagName) != atom.Span { goto nextToken }
//
//			state = 0
//
//		// EOF
//		case t == html.ErrorToken:
//			goto doneParsing
//		}
//
//	nextToken:
//	}
//doneParsing:
//
//	return
//}
//
//func validText(b []byte) bool {
//	if len(b) == 0 { return false }
//	if b[0] == ' ' || b[0] == '\n' || b[0] == ',' {
//		return false
//	}
//	return true
//}
