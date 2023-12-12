package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/smekuria1/og/sitemap"
	"golang.org/x/sync/errgroup"
)

var urlstr = flag.String("urlstr", "", "url to build sitemap from")
var depth = flag.Int("depth", 1, "depth to crawl")
var outFile = flag.String("out", "./map.xml", "output file path")
var inFile = flag.String("in", "", "input file of urls to buildsitemap")

func main() {
	flag.Parse()
	if *inFile != "" {
		g := new(errgroup.Group)
		urls := readlines(*inFile)
		fmt.Printf("urls: %v\n", urls)
		for i, u := range urls {
			url := u
			i := i
			g.Go(func() error {
				return sitemap.Sitemap(url, *depth, fmt.Sprintf("./map%v.xml", i))
			})
		}
		if err := g.Wait(); err == nil {
			fmt.Println("Finished building Sitemaps")
		}

		return

	} else if *urlstr == "" {
		log.Fatal("url is required")
	}

	if *depth < 0 {
		log.Fatal("depth must be greater than 0")
	}
	sitemap.Sitemap(*urlstr, *depth, *outFile)
}

func readlines(input string) []string {
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	file.Close()

	return lines
}
