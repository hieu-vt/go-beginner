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

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func CrawlRecursive(url string, depth int, fetcher Fetcher, quit chan bool, visitedUrls map[string]bool) {
	if depth <= 0 {
		quit <- true
		return
	}

	// didIt, hasIt := visitedUrls[url]
	// If we have already visited this link,
	// stop here
	// if didIt && hasIt {
	// 	quit <- true
	// 	return

	// } else {
	// Mark it has visited
	visitedUrls[url] = true
	// }

	// Fetch children URLs
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println("Error: ", url)

		quit <- true
		return
	}

	fmt.Println(body)

	// Crawl children URLs
	childrenQuit := make(chan bool)
	for _, childrenUrl := range urls {
		go CrawlRecursive(childrenUrl, depth-1, fetcher, childrenQuit, visitedUrls)
		// To exit goroutines. This channel will always be filled
		<-childrenQuit
	}

	quit <- true
}

func Crawl(url string, depth int, fetcher Fetcher) {
	quit := make(chan bool)
	// Say we haven't visited the first URL yet
	visitedUrls := map[string]bool{url: false}

	// Le'ts go, crawl from the given URL
	go CrawlRecursive(url, depth, fetcher, quit, visitedUrls)

	// We will not quit until we have something
	// in the "quit" channel
	<-quit
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// Fetch returns the body of URL and
// a slice of URLs found on that page.
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
