package wait_test

import (
	"context"
	"fmt"
	"time"

	"github.com/jaslong/wait"
)

func Example_context() {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	ctxDoneSignal := wait.For(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	taskSignal := wait.For(func() error {
		time.Sleep(2 * time.Second)
		return nil
	})

	err := wait.UntilAllSucceed(ctxDoneSignal, taskSignal).Wait()
	if err == nil {
		fmt.Println("this should not be printed")
	}

	fmt.Println(err.Error())
	// Output: context deadline exceeded
}
