package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

var concurrency = pflag.UintP("concurrency", "c", 4, "Number of concurrent downloads")
var dir = pflag.StringP("out-dir", "o", "Downloads", "Output directory")
var verbose = pflag.BoolP("verbose", "v", false, "More output")

var availableLetters = [...]string{
	//"#",
	//"A",
	//"B",
	//"C",
	//"D",
	//"E",
	//"F",
	//"G",
	//"H",
	//"I",
	//"J",
	//"K",
	//"L",
	//"M",
	//"N",
	//"O",
	//"P",
	"Q",
	//"R",
	//"S",
	//"T",
	//"U",
	//"V",
	//"W",
	//"X",
	//"Y",
	//"Z",
}

func parseArgs() error {
	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr,
			`
Usage:`)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if *concurrency <= 0 {
		return fmt.Errorf("invalid value for --concurrency")
	}

	if err := os.MkdirAll(*dir, 0777); err != nil {
		return err
	}

	return nil
}
