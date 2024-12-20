package polipo

import (
	"context"
	"errors"
	"sync"
)

// defaultMaxConcurrency is the default maximum number of concurrent tasks to run.
const defaultMaxConcurrency = 10

// Task is a function that can be run by the Polipo.
// Arguments could be passed in as a closure.
// The function should return a slice of generic items and an error.
type Task[T any] func() ([]T, error)

// Polipo stores a list of Tasks to be run concurrently.
type Polipo[T any] struct {
	tasks             []Task[T]
	maxConcurrency    int
	concurrencyBuffer chan struct{}
}

// NewPolipo creates a new Polipo. It accepts options to configure the Polipo.
// The default maximum number of concurrent tasks is 10.
// The options are:
// - WithMaxConcurrency: sets the maximum number of concurrent tasks to run.
func NewPolipo[T any](opts ...Option[T]) Polipo[T] {
	p := Polipo[T]{
		tasks:          make([]Task[T], 0),
		maxConcurrency: defaultMaxConcurrency,
	}

	for _, opt := range opts {
		opt(&p)
	}

	p.concurrencyBuffer = make(chan struct{}, p.maxConcurrency)

	// Fill the concurrencyBuffer with available slots.
	for i := 0; i < p.maxConcurrency; i++ {
		p.concurrencyBuffer <- struct{}{}
	}

	return p
}

// AddTask adds a Task to the Polipo. The Task function will be run when Do is called.
func (p *Polipo[T]) AddTask(task Task[T]) {
	p.tasks = append(p.tasks, task)
}

// Do executes all the Tasks concurrently. It limits the number of concurrent tasks to the value
// set by passing `WithMaxConcurrency` as an option. The default is 10.
// This is a blocking function that will return when all the Tasks have finished their work.
func (p *Polipo[T]) Do(ctx context.Context) ([]T, error) {
	if len(p.tasks) == 0 {
		return nil, errors.New("no tasks to do")
	}

	resultsChan := make(chan result[T])
	wg := sync.WaitGroup{}

	wg.Add(len(p.tasks))

	// Schedule tasks to run concurrently limiting the number of concurrent tasks.
	go func() {
		for _, task := range p.tasks {
			// Wait for an available slot in the concurrencyBuffer.
			<-p.concurrencyBuffer

			go func(t Task[T]) {
				defer wg.Done()
				items, err := t()

				select {
				case resultsChan <- result[T]{items, err}:
				case <-ctx.Done():
				}

				// Release the slot in the concurrencyBuffer to allow other tasks to run.
				p.concurrencyBuffer <- struct{}{}
			}(task)
		}
	}()

	// Wait for all tasks to finish.
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var (
		results []T
		errs    []error
	)

	// Collect results and errors.
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		select {
		case res, ok := <-resultsChan:
			if !ok {
				return results, errors.Join(errs...)
			}

			if res.err != nil {
				errs = append(errs, res.err)
			}

			results = append(results, res.items...)
		case <-ctx.Done():
			return results, errors.Join(append(errs, ctx.Err())...)
		}
	}
}

type result[T any] struct {
	items []T
	err   error
}
