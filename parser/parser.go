package parser

import (
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func textNodes(n *html.Node) (ret string) {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += textNodes(c) + " "
	}
	return strings.Join(strings.Fields(ret), " ")
}

func buildLink(n *html.Node) Link {
	return Link{Href: n.Attr[0].Val, Text: textNodes(n)}
}

func Parse(htmlCont []byte) (linkArr []Link, err error) {
	doc, err := html.Parse(strings.NewReader(string(htmlCont)))
	if err != nil {
		return nil, err
	}

	linkNodes := linkNodes(doc)
	for _, node := range linkNodes {
		linkArr = append(linkArr, buildLink(node))
	}
	return linkArr, nil
}
