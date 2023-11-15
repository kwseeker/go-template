package sync

import (
	"sync"
	"testing"
)

type httpPkg struct{}

func (httpPkg) Get(url string) {}

var http httpPkg

// This example fetches several URLs concurrently,
// using a WaitGroup to block until all the fetches are complete.
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			http.Get(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}
