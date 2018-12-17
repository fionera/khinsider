package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
)

func download(u string, t Track) (int64, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(u)

	if err := fasthttp.Do(req, res); err != nil {
		return 0, err
	}

	if sc := res.StatusCode(); sc != 200 {
		return 0, fmt.Errorf("HTTP status %d", sc)
	}

	folder := filepath.Join(*dir, string(t.Game.Name))

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return 0, err
	}

	fileName := filepath.Join(folder, filepath.Base(u))

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		logrus.Infof("Skipping File " + fileName)
		return 0, nil
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return 0, err
	}

	size, err := res.WriteTo(f)
	if err != nil {
		return 0, err
	}

	return size, f.Close()
}
