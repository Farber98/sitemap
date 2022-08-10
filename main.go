package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
	1. Get webpage.
	2. Parse all links on the page.
	3. Build proper urls with our links
	4. Filter out any links w/ different domains.
	5. Find all pages (BFS) and repeat 1 to 4.
	6. Print out XML
*/

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	flag.Parse()

	fmt.Println(*urlFlag)

	html, err := getWebpage(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", html)
}

func getWebpage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body))
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}
