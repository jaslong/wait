// Package wait is a simple yet powerful Go concurrency abstraction.
//
// In Go, a goroutine starts an operation asynchronously but doesn't provide mechanisms for getting
// notified of its completion. Instead, Go provides channels for asynchronous message passing and
// requires you to use them for each goroutine.
//
// Let's say we have a func DoWork() (Result, error) function. We want to run it asynchronously.
// We would normally use channels to communicate the result in the goroutine to the main function:
//
//  func main() {
//    var result Result
//    errChan := make(chan error)
//    go func() {
//      var err error
//      result, err = DoWork()
//      errChan <- err
//    }()
//
//    err := <-errChan
//    if err != nil {
//      fmt.Printf("Failure: %v", err)
//    } else {
//      fmt.Printf("Success: %v" + result)
//    }
//  }
//
// With wait.For, you don't need to worry about channels:
//
//  func main() {
//    var result Result
//    signal := wait.For(func() (err error) {
//      result, err = doWork()
//      return
//    })
//
//    err := signal.Wait()
//    if err != nil {
//      fmt.Printf("Failure: %v", err)
//    } else {
//      fmt.Printf("Success: %v", result)
//    }
//  }
//
// wait.For takes a func() error, runs it in a goroutine, and returns a Signal. A Signal is
// a simple interface with a single method Wait, which blocks until the func() error returns.
//
// Signals can be easily merged in many ways. This package comes with a few:
//
// * wait.UntilAllComplete waits until all input Signals complete, regardless of whether they
// succeed or fail. Equivalent to sync.WaitGroup: https://godoc.org/sync#WaitGroup.
//
// * wait.UntilAllSucceed waits until all input Signals succeed, failing immediately if any
// Signal fails. Equivalent to errgroup.Group: https://godoc.org/golang.org/x/sync/errgroup.
//
// Let's say we have another function func DoOtherWork() (OtherResult, error) function and want to
// run it concurrently with DoWork. With normal Go:
//
//  func main() {
//    var (
//      result Result
//      otherResult OtherResult
//    )
//    errChan := make(chan error)
//    go func() {
//      var err error
//      result, err = DoWork()
//      errChan <- err
//    }()
//    go func() {
//      var err error
//      otherResult, err = DoOtherWork()
//      errChan <- err
//    }()
//
//    for i := 0; i < 2; i++ {
//      err := <-errChan
//      if err != nil {
//        fmt.Printf("Failure: %v", err)
//        return
//      }
//    }
//    fmt.Printf("Success: %v, %v", result, otherResult)
//  }
//
// With wait:
//
//  func main() {
//    var (
//      result Result
//      otherResult OtherResult
//    )
//    signal := wait.For(func() (err error) {
//      result, err = DoWork()
//      return
//    })
//    otherSignal := wait.For(func() (err error) {
//      otherResult, err = DoOtherWork()
//      return
//    })
//
//    mergedSignal := wait.UntilAllSucceed(signal, otherSignal)
//    err := mergedSignal.Wait()
//    if err != nil {
//      fmt.Printf("Failure: %v", err)
//    } else {
//      fmt.Printf("Success: %v, %s", result, otherResult)
//    }
//  }
package wait

// Signal represents completion of an abstract operation.
type Signal interface {
	// Wait blocks the underlying operation completes.
	Wait() error
}

// For starts a goroutine that runs the routine and returns a signal representing its completion.
func For(routine func() error) Signal {
	signal := channelSignal{
		channel: make(chan struct{}),
	}
	go func() {
		defer close(signal.channel)
		signal.err = routine()
		signal.done = true
	}()
	return &signal
}

type channelSignal struct {
	channel chan struct{}
	done    bool
	err     error
}

func (s *channelSignal) Wait() error {
	if !s.done {
		<-s.channel
	}
	return s.err
}
