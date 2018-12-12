package main

import "context"

type Track struct {
	Game     Game
	Title    string
	Download string
}

func (t *Track) Crawl(c context.Context) error {
	return nil
}