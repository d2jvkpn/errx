package errx

import (
	"errors"
	//"fmt"
	"sync"
)

func ParallelRun(funcs ...func() *ErrX) (err *ErrX) {
	var (
		errs []*ErrX
		wg   sync.WaitGroup
	)

	errs = make([]*ErrX, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := range errs {
		if errs[i] == nil {
			continue
		}

		if err == nil {
			err = errs[i]
		} else {
			err.WithErr(errs[i])
		}
	}

	return err
}

func ParallelRun2(funcs ...func() error) (err *ErrX) {
	var (
		hasError bool
		ok       bool
		errs     []error
		wg       sync.WaitGroup
	)

	errs = make([]error, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}
	wg.Wait()

	hasError = false
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		hasError = true

		if err != nil {
			if err, ok = errs[i].(*ErrX); ok {
				errs[i] = nil
			}
		}
	}

	if !hasError {
		return nil
	}

	if err == nil {
		err = NewErrX(errors.Join(errs...))
	} else {
		err.WithErr(errs...)
	}

	return err
}
