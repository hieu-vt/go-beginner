/**
Đề bài: Cho 10000 URLs
Tạo 5 routines để nó có thể crawl data 1 lúc từ 5 url về
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func createUrls(num int) []int {
	urls := make([]int, num)

	for i := 0; i < num; i++ {
		urls[i] = i + rand.Intn(10*(i+1))
	}

	return urls
}

func crawlData(url int, name string) {
	time.Sleep(time.Millisecond * time.Duration(url))
	fmt.Printf("Worker %s run %d \n", name, url)
}

func main() {
	const maxRoutine = 5
	const maxRequest = 100
	dataUrls := createUrls(maxRequest)
	quoteChan := make(chan int, maxRequest)

	for i := 0; i < len(dataUrls); i++ {
		quoteChan <- dataUrls[i]
	}

	for i := 1; i <= maxRoutine; i++ {
		go func(name string) {
			for value := range quoteChan {
				crawlData(value, name)
			}

			fmt.Printf("Close worker %s", name)
		}(fmt.Sprintf("%d", i))
	}

	close(quoteChan)
	// time.Sleep(time.Second * 10)
}
