package wait_test

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/jaslong/wait"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UntilAllComplete", func() {
	for n := 0; n <= 10; n++ {
		It(fmt.Sprintf("waits until all %d routines are done", n), func() {
			signals := make([]wait.Signal, n)
			returnFuncs := make([]func(error), n)
			for i := 0; i < n; i++ {
				signal, returnFunc := waitForTestRoutine()
				signals[i] = signal
				returnFuncs[i] = returnFunc
			}

			errChan := ErrChan(wait.UntilAllComplete(signals...))

			// shuffle the order routines return
			rand.Shuffle(n, func(i, j int) {
				returnFuncs[i], returnFuncs[j] = returnFuncs[j], returnFuncs[i]
			})

			for i := 0; i < n; i++ {
				// whether the routine errors or not is arbitrary, so make it random
				var err error
				if rand.Int()%2 == 0 {
					err = errors.New("!")
				}
				returnFuncs[i](err)

				if i < n-1 {
					Expect(errChan).To(BeEmpty())
					Expect(errChan).NotTo(BeClosed())
				}
			}
			Expect(<-errChan).To(BeNil())
			Eventually(errChan).Should(BeClosed())
		})
	}
})

var _ = Describe("UntilAllSucceed", func() {
	for n := 0; n <= 10; n++ {
		It(fmt.Sprintf("waits until all %d routines succeed", n), func() {
			signals := make([]wait.Signal, n)
			returnFuncs := make([]func(error), n)
			for i := 0; i < n; i++ {
				signal, returnFunc := waitForTestRoutine()
				signals[i] = signal
				returnFuncs[i] = returnFunc
			}

			errChan := ErrChan(wait.UntilAllSucceed(signals...))

			// shuffle the order routines return
			rand.Shuffle(n, func(i, j int) {
				returnFuncs[i], returnFuncs[j] = returnFuncs[j], returnFuncs[i]
			})

			for i := 0; i < n; i++ {
				returnFuncs[i](nil)

				if i < n-1 {
					Expect(errChan).To(BeEmpty())
					Expect(errChan).NotTo(BeClosed())
				}
			}
			Expect(<-errChan).To(BeNil())
			Eventually(errChan).Should(BeClosed())
		})
	}

	for n := 0; n <= 10; n++ {
		It(fmt.Sprintf("waits until all %d routines are done, 1 fails", n), func() {
			// make 1 routine fail
			fixedErr := errors.New("!")
			failIndex := rand.Intn(n)

			signals := make([]wait.Signal, n)
			returnFuncs := make([]func(error), n)
			for i := 0; i < n; i++ {
				signal, returnFunc := waitForTestRoutine()
				signals[i] = signal
				returnFuncs[i] = returnFunc
			}

			errChan := ErrChan(wait.UntilAllSucceed(signals...))

			// shuffle the order routines return
			rand.Shuffle(n, func(i, j int) {
				returnFuncs[i], returnFuncs[j] = returnFuncs[j], returnFuncs[i]
			})

			var hasFailed bool
			for i := 0; i < n; i++ {
				if i == failIndex {
					hasFailed = true
					returnFuncs[i](fixedErr)
				} else {
					returnFuncs[i](nil)
				}

				// we can only expect errChan to be empty and not closed before failing
				if !hasFailed && i < n-1 {
					Expect(errChan).To(BeEmpty())
					Expect(errChan).NotTo(BeClosed())
				}
			}
			Expect(<-errChan).To(BeIdenticalTo(fixedErr))
			Eventually(errChan).Should(BeClosed())
		})
	}
})
