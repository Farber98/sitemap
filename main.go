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
	flag.Parse()

	finalLinks, err := getWebpage(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}
	for _, link := range finalLinks {
		if _, ok := linkMap[link]; !ok {
			linkMap[link] = link
		}
	}

	for k := range linkMap {
		fmt.Println(k)
	}
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
