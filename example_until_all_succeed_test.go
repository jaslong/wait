package wait_test

import (
	"context"
	"fmt"
	"os"

	"github.com/jaslong/wait"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

// This example is equivalent to the errgroup.Group parallel example found at
// https://godoc.org/golang.org/x/sync/errgroup#example-Group--Parallel.
func ExampleUntilAllSucceed() {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		searches := []Search{Web, Image, Video}
		signals := make([]wait.Signal, len(searches))
		results := make([]Result, len(searches))
		for i, search := range searches {
			i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
			signals[i] = wait.For(func() error {
				result, err := search(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := wait.UntilAllSucceed(signals...).Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
