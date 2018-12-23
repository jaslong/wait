package wait_test

import (
	"testing"

	"github.com/jaslong/wait"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWait(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wait")
}

// ErrChan wraps the signal and returns a one-time-use channel that receives the error.
//
// TODO: Consider exporting this function.
func ErrChan(signal wait.Signal) <-chan error {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		errChan <- signal.Wait()
	}()
	return errChan
}
