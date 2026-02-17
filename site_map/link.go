package sitemap

import (
	"io"

	"golang.org/x/net/html"
)

func Parse(r io.Reader) ([]string, error) {

	node, _ := html.Parse(r)

	var links []string

	extractLinks(node, &links)

	return links, nil
}

func extractLinks(node *html.Node, links *[]string) {

	if node.Type == html.ElementNode && node.Data == "a" {
		var link string
		for _, a := range node.Attr {
			if a.Key == "href" {
				link = a.Val
				break
			}
		}
		*links = append(*links, link)
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		extractLinks(n, links)
	}
}
