package wait_test

import (
	"net/http"

	"github.com/jaslong/wait"
)

// This example is equivalent to the sync.WaitGroup example found at
// https://godoc.org/sync#example-WaitGroup.
func ExampleUntilAllComplete() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}

	signals := make([]wait.Signal, len(urls))
	for i, url := range urls {
		// Copy url to a local variable since url will be mutated during iteration.
		localURL := url
		// Launch a goroutine to fetch the URL.
		signals[i] = wait.For(func() error {
			// Fetch the URL.
			http.Get(localURL)
			return nil
		})
	}
	// Wait for all HTTP fetches to complete.
	wait.UntilAllComplete(signals...).Wait()
}
