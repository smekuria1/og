/*
Package sitemap provides a simple way to build a sitemap for a given url.
1. GET the Web Page
2. Parse all the links on the page
3. Build proper urls from the links
4. Filter out any links with a different domain
5. Find all the pages (BFS)
6. Output the data as XML
*/
package sitemap

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/smekuria1/og/link"
)

// Sitemap function
func Sitemap(urlstr string, depth int, out string) {
	log.Printf("Writing To %s", out)
	pages := bfs(urlstr, depth)
	toXML := urlset{
		Xmlns: xmlns,
	}

	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}
	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(xml.Header)

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")

	if encErr := enc.Encode(toXML); encErr != nil {
		panic(encErr)
	}
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func bfs(urlstr string, depth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	newQ := map[string]struct{}{
		urlstr: {},
	}
	for i := 0; i <= depth; i++ {
		q, newQ = newQ, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for urlstr := range q {
			if _, ok := seen[urlstr]; ok {
				continue
			}
			seen[urlstr] = struct{}{}
			for _, link := range get(urlstr) {
				newQ[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for urlstr := range seen {
		ret = append(ret, urlstr)
	}
	return ret

}

func filter(keepFn func(string) bool, links []string) []string {
	var filtered []string
	for _, link := range links {
		if keepFn(link) {
			filtered = append(filtered, link)
		}
	}
	return filtered
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}

func get(urlstr string) []string {

	resp, err := http.Get(urlstr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	baseUrl := &url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}
	base := baseUrl.String()

	return filter(withPrefix(base), hrefs(resp.Body, base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)

	var hrefs []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			hrefs = append(hrefs, base+link.Href)
		case strings.HasPrefix(link.Href, "http"):
			hrefs = append(hrefs, link.Href)
			// default:
			// 	log.Println("Skipping", link.Href)
		}
	}
	return hrefs
}
