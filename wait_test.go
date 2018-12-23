package wait_test

import (
	"errors"

	"github.com/jaslong/wait"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("For", func() {
	It("waits for a routine to return", func() {
		signal, returnFunc := waitForTestRoutine()
		errChan := ErrChan(signal)

		Expect(errChan).To(BeEmpty())
		Expect(errChan).NotTo(BeClosed())

		returnFunc(nil)
		Expect(<-errChan).To(BeNil())
		Eventually(errChan).Should(BeClosed())
	})

	It("waits for a routine to return without error", func() {
		fixedErr := errors.New("!")

		signal, returnFunc := waitForTestRoutine()
		errChan := ErrChan(signal)

		Expect(errChan).To(BeEmpty())
		Expect(errChan).NotTo(BeClosed())

		returnFunc(fixedErr)
		Expect(<-errChan).To(BeIdenticalTo(fixedErr))
		Eventually(errChan).Should(BeClosed())
	})
})

func waitForTestRoutine() (wait.Signal, func(error)) {
	returnFunc := make(chan error)
	return wait.For(func() error {
			return <-returnFunc
		}),
		func(err error) {
			returnFunc <- err
		}
}
