package parallel

import "sync"

type Result[I any, O any] struct {
	Input  I
	Output O
	Err    error
}

type WorkFunc[I any, O any] func(input I) (O, error)

func Do[I any, O any](work []I, workFunc WorkFunc[I, O], workers int, bufferSize int) []Result[I, O] {
	workC := make(chan I, bufferSize)
	resultC := make(chan Result[I, O], bufferSize)
	results := make([]Result[I, O], 0, len(work))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for w := range workC {
				o, err := workFunc(w)
				resultC <- Result[I, O]{w, o, err}
			}
		}()
	}

	// load up the work
	go func() {
		defer close(workC)
		for _, w := range work {
			workC <- w
		}
	}()

	// close results channel if all workers are done
	go func() {
		wg.Wait()
		close(resultC)
	}()

	// capture all results
	for r := range resultC {
		results = append(results, r)
	}

	return results
}

type StateFunc[S any] func() S

type WorkFuncWithState[S any, I any, O any] func(state S, input I) (O, error)

func DoWithState[I any, O any, S any](work []I, stateFunc StateFunc[S], workFunc WorkFuncWithState[S, I, O], workers int, bufferSize int) []Result[I, O] {
	workC := make(chan I, bufferSize)
	resultC := make(chan Result[I, O], bufferSize)
	results := make([]Result[I, O], 0, len(work))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			s := stateFunc()
			defer wg.Done()
			for w := range workC {
				o, err := workFunc(s, w)
				resultC <- Result[I, O]{w, o, err}
			}
		}()
	}

	// load up the work
	go func() {
		defer close(workC)
		for _, w := range work {
			workC <- w
		}
	}()

	// close results channel if all workers are done
	go func() {
		wg.Wait()
		close(resultC)
	}()

	// capture all results
	for r := range resultC {
		results = append(results, r)
	}

	return results
}
