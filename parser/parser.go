package parser

import (
	"strings"

	"golang.org/x/net/html"
)

var visited = map[string]bool{
	"/": true,
}

type link struct {
	href string
	text string
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

func buildLink(n *html.Node) link {
	visited[n.Attr[0].Val] = true
	return link{href: n.Attr[0].Val, text: textNodes(n)}
}

func Parse(htmlCont []byte, domain string) (linkArr []link, err error) {
	doc, err := html.Parse(strings.NewReader(string(htmlCont)))
	if err != nil {
		return nil, err
	}

	linkNodes := linkNodes(doc)
	for _, node := range linkNodes {
		if !shouldParse(node) {
			continue
		}
		node = normalizeSameDomain(node, domain)
		if diffDomain(node, domain) {
			continue
		}
		linkArr = append(linkArr, buildLink(node))
	}
	return linkArr, nil
}

func shouldParse(n *html.Node) bool {
	text := n.Attr[0].Val

	//not a link
	if !strings.HasPrefix(text, "http") && !strings.HasPrefix(text, "/") {
		return false
	}

	// already visited
	if _, ok := visited[n.Attr[0].Val]; ok {
		return false
	}

	return true
}

func normalizeSameDomain(n *html.Node, domain string) *html.Node {
	if strings.HasPrefix(n.Attr[0].Val, "/") {
		domain += n.Attr[0].Val
		n.Attr[0].Val = domain
	}
	return n
}

func diffDomain(n *html.Node, domain string) bool {
	if !strings.HasPrefix(n.Attr[0].Val, domain) {
		return true
	}
	return false
}
