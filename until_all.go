package wait

// UntilAllComplete returns a new signal that succeeds when all input signals complete, regardless
// of success/failure.
func UntilAllComplete(signals ...Signal) Signal {
	return &allDoneSignal{
		signals: signals,
	}
}

type allDoneSignal struct {
	signals []Signal
}

func (f *allDoneSignal) Wait() error {
	for _, signal := range f.signals {
		_ = signal.Wait() // nolint: gosec
	}
	return nil
}

// UntilAllSucceed returns a new signal that succeeds if all input signals succeed. If any input
// signal fails, the returned signal fails with the error of one of the failed signals.
func UntilAllSucceed(signals ...Signal) Signal {
	return &allSucceedSignal{
		signals: signals,
	}
}

type allSucceedSignal struct {
	signals []Signal
}

func (f *allSucceedSignal) Wait() error {
	errChan := make(chan error)
	for _, signal := range f.signals {
		go func(signal Signal) {
			errChan <- signal.Wait()
		}(signal)
	}
	for i := 0; i < len(f.signals); i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}
