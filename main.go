package main

import (
	"fmt"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, channel chan string) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:

	defer close(channel)

	if depth <= 0 {

		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		channel <- "Error: " + url

		return
	}

	channel <- body

	chanUrls := make([]chan string, len(urls))

	for i, u := range urls {
		chanUrls[i] = make(chan string)
		go Crawl(u, depth-1, fetcher, chanUrls[i])
	}

	for j := range chanUrls {
		for c := range chanUrls[j] {
			channel <- c
		}
	}

	return
}

func main() {
	channel := make(chan string)

	go Crawl("https://golang.org/", 4, fetcher, channel)

	for c := range channel {
		fmt.Println(c)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if url == "https://golang.org/pkg/os/" {
		time.Sleep(time.Microsecond * 200)
	}

	if url == "https://golang.org/" {
		time.Sleep(time.Millisecond * 300)
	}

	time.Sleep(time.Millisecond * 100)

	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"111",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"222",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"333",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"444",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
