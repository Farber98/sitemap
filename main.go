package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sitemap-builder/parser"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

/*
	1. Get webpage.
	2. Parse all links on the page.
		2.1 Build proper urls with our links
		2.2 Filter out any links w/ different domains.
	3. Print out XML
*/

func main() {
	urlFlag := flag.String("url", "https://github.com/Farber98?tab=repositories", "the url that you want to build a sitemap for")
	flag.Parse()

	fmt.Println(*urlFlag)

	html, err := getWebpage(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	links, err := parser.Parse(html, *urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	sm := buildXML(links, *urlFlag)

	fmt.Println(string(sm.XMLContent()))
}

func getWebpage(url string) ([]byte, error) {
	res, err := http.Get(url)
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
	return body, nil
}

func buildXML(links []parser.Link, domain string) *stm.Sitemap {
	sm := stm.NewSitemap()
	sm.SetDefaultHost(domain)

	sm.Create()

	for _, link := range links {
		sm.Add(stm.URL{"loc": link.Href})
	}
	return sm
}
