# wait

**`wait` is a simple yet powerful Go concurrency abstraction**

In Go, a *goroutine* starts an operation asynchronously but doesn't provide mechanisms for getting
notified of its completion. Instead, Go provides *channels* for asynchronous message passing and
requires you to use them for each goroutine.

Let's say we have a `func DoWork() (Result, error)` function. We want to run it asynchronously.
We would normally use channels to communicate the result in the goroutine to the main function:

```go
func main() {
  var result Result
  errChan := make(chan error)
  go func() {
    var err error
    result, err = DoWork()
    errChan <- err
  }()

  err := <-errChan
  if err != nil {
    fmt.Printf("Failure: %v", err)
  } else {
    fmt.Printf("Success: %v" + result)
  }
}
```

With `wait.For`, you don't need to worry about channels:

```go
func main() {
  var result Result
  signal := wait.For(func() (err error) {
    result, err = doWork()
    return
  })

  err := signal.Wait()
  if err != nil {
    fmt.Printf("Failure: %v", err)
  } else {
    fmt.Printf("Success: %v", result)
  }
}
```

`wait.For` takes a `func() error`, runs it in a goroutine, and returns a `Signal`. A `Signal` is
a simple interface with a single method `Wait`, which blocks until the `func() error` returns.

`Signal`s can be easily merged in many ways. This package comes with a few:

* `wait.UntilAllComplete` waits until all input `Signal`s complete, regardless of whether they
succeed or fail. Equivalent to [`sync.WaitGroup`](https://godoc.org/sync#WaitGroup).
* `wait.UntilAllSucceed` waits until all input `Signal`s succeed, failing immediately if any
`Signal` fails. Equivalent to [`errgroup.Group`](https://godoc.org/golang.org/x/sync/errgroup).

Let's say we have another function `func DoOtherWork() (OtherResult, error)` function and want to
run it concurrently with `DoWork`. With normal Go:

```go
func main() {
  var (
    result Result
    otherResult OtherResult
  )
  errChan := make(chan error)
  go func() {
    var err error
    result, err = DoWork()
    errChan <- err
  }()
  go func() {
    var err error
    otherResult, err = DoOtherWork()
    errChan <- err
  }()

  for i := 0; i < 2; i++ {
    err := <-errChan
    if err != nil {
      fmt.Printf("Failure: %v", err)
      return
    }
  }
  fmt.Printf("Success: %v, %v", result, otherResult)
}
```

With `wait`:

```go
func main() {
  var (
    result Result
    otherResult OtherResult
  )
  signal := wait.For(func() (err error) {
    result, err = DoWork()
    return
  })
  otherSignal := wait.For(func() (err error) {
    otherResult, err = DoOtherWork()
    return
  })

  mergedSignal := wait.UntilAllSucceed(signal, otherSignal)
  err := mergedSignal.Wait()
  if err != nil {
    fmt.Printf("Failure: %v", err)
  } else {
    fmt.Printf("Success: %v, %v", result, otherResult)
  }
}
```
