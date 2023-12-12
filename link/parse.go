package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an HTML link (<a href="...">) in an HTML
// document.
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return
// a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
		//fmt.Println(links)
	}
	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var out []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		out = append(out, linkNodes(c)...)
	}
	return out
}

func buildLink(n *html.Node) Link {
	var out Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			out.Href = attr.Val
			break
		}
	}
	out.Text = buildText(n)
	return out
}

func buildText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var builder strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		builder.WriteString(buildText(c) + " ")
	}
	return strings.Join(strings.Fields(builder.String()), " ")
}
