package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sitemap-builder/parser"
	"strings"
)

/*
	1. Get webpage.
	2. Parse all links on the page.
		2.1 Build proper urls with our links
		2.2 Filter out any links w/ different domains.
	3. Print out XML
*/

var linkMap = map[string]string{}

func main() {
	urlFlag := flag.String("url", "https://epipe.com/", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "maximum number of links deep to traverse")
	flag.Parse()

	links, err := bfs(*urlFlag, *maxDepth)
	if err != nil {
		log.Fatal(err)
	}
	for _, link := range links {
		fmt.Println(link)
	}
}

func bfs(urlStr string, maxDepth int) (ret []string, err error) {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{urlStr: {}}
	for i := 0; i < maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			links, err := getWebpage(url)
			if err != nil {
				return nil, err
			}
			for _, link := range links {
				nq[link] = struct{}{}
			}
		}
	}
	for url := range seen {
		ret = append(ret, url)
	}
	return ret, nil
}

func cleanLinks(links []parser.Link, domain string) (hrefs []string) {
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			hrefs = append(hrefs, domain+link.Href)
		case strings.HasPrefix(link.Href, "http"):
			hrefs = append(hrefs, link.Href)
		default:
		}
	}
	return
}

func getWebpage(domain string) ([]string, error) {
	res, err := http.Get(domain)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body))
	}
	if err != nil {
		return nil, err
	}

	baseUrl := &url.URL{
		Scheme: res.Request.URL.Scheme,
		Host:   res.Request.URL.Host,
	}

	links, _ := parser.Parse(body)

	cleanLinks := cleanLinks(links, baseUrl.String())

	return filter(baseUrl.String(), cleanLinks), nil
}

func filter(base string, links []string) (filteredLinks []string) {
	for _, link := range links {
		if strings.HasPrefix(link, base) {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return
}
