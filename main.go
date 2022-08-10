package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sitemap-builder/parser"

	"github.com/snabb/sitemap"
)

/*
	1. Get webpage.
	2. Parse all links on the page.
		2.1 Build proper urls with our links
		2.2 Filter out any links w/ different domains.
	3. Print out XML
*/

func main() {
	urlFlag := flag.String("url", "https://www.calhoun.io", "the url that you want to build a sitemap for")
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

	sm.WriteTo(os.Stdout)
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

func buildXML(links []parser.Link, domain string) *sitemap.Sitemap {
	sm := sitemap.New()
	for _, link := range links {
		sm.Add(&sitemap.URL{Loc: link.Href})
	}
	return sm
}
