package wait_test

import (
	"errors"
	"fmt"

	"github.com/jaslong/wait"
)

func ExampleFor_success() {
	signal := wait.For(func() error {
		fmt.Println("routine succeeding")
		return nil
	})

	err := signal.Wait()
	if err != nil {
		fmt.Println("this should not happen")
	}
	fmt.Println("routine succeeded")

	// Output:
	// routine succeeding
	// routine succeeded
}

func ExampleFor_failure() {
	signal := wait.For(func() error {
		fmt.Println("routine failing")
		return errors.New("routine failed")
	})

	err := signal.Wait()
	if err == nil {
		fmt.Println("this should not happen")
	}
	fmt.Println(err.Error())

	// Output:
	// routine failing
	// routine failed
}
