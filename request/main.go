/**
Đề bài: Cho 10000 URLs
Tạo 5 routines để nó có thể crawl data 1 lúc từ 5 url về
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Storage struct {
	count int
	store map[string]string
	lock  *sync.RWMutex
}

func (s *Storage) Read(k string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store[k]
}

func (s *Storage) Write(k string, v string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.store[k] = v
	s.count++
}

func NewStorage() *Storage {
	return &Storage{
		count: 0,
		store: make(map[string]string),
		lock:  new(sync.RWMutex),
	}
}

func createUrls(num int) []int {
	urls := make([]int, num)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < num; i++ {
		urls[i] = r1.Intn(10)
	}

	return urls
}

func crawlData(url int, name string) {
	time.Sleep(time.Second / 10 * time.Duration(url))
	fmt.Printf("Worker %s run %d \n", name, url)
}

func main() {
	const maxRoutine = 5
	const maxRequest = 100
	dataUrls := createUrls(maxRequest)
	doneChan := make(chan int)
	defer close(doneChan)
	quoteChan := make(chan int, maxRequest)
	store := NewStorage()

	for i := 0; i < len(dataUrls); i++ {
		quoteChan <- dataUrls[i]
	}

	for i := 1; i <= maxRoutine; i++ {
		go func(name string) {
			for value := range quoteChan {
				key := fmt.Sprintf("%d", value)
				visitedUrl := store.Read(key)
				if visitedUrl != "" {
					continue
				}
				crawlData(value, name)
				store.Write(key, key)
			}

			fmt.Printf("Close worker %s\n", name)
			doneChan <- 1
		}(fmt.Sprintf("%d", i))
	}

	close(quoteChan)

	for i := 1; i <= maxRoutine; i++ {
		<-doneChan
	}

	fmt.Println("Visited url: ", store.count)
}
