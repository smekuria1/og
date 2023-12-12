package main

import (
	"flag"
	"log"

	"github.com/smekuria1/og/sitemap"
)

var urlstr = flag.String("urlstr", "", "url to build sitemap from")
var depth = flag.Int("depth", 1, "depth to crawl")
var outFile = flag.String("out", "./map.xml", "output file path")

func main() {
	flag.Parse()
	if *urlstr == "" {
		log.Fatal("url is required")
	}
	if *depth < 0 {
		log.Fatal("depth must be greater than 0")
	}
	sitemap.Sitemap(*urlstr, *depth, *outFile)
}
