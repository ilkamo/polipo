package polipo

import (
	"context"
	"errors"
	"sync"
)

// Task is a function that can be run by the Polipo.
// Arguments could be passed in as a closure.
// The function should return a slice of generic items and an error.
type Task[T any] func() ([]T, error)

// Polipo stores a list of Tasks to be run concurrently.
type Polipo[T any] struct {
	tasks []Task[T]
}

// NewPolipo creates a new Polipo.
func NewPolipo[T any]() Polipo[T] {
	return Polipo[T]{
		tasks: make([]Task[T], 0),
	}
}

// AddTask adds a Task to the Polipo. The Task function will be run when Do is called.
func (p *Polipo[T]) AddTask(task Task[T]) {
	p.tasks = append(p.tasks, task)
}

// Do executes all the Tasks concurrently.
// This is a blocking function that will return when all the Tasks have finished their work.
func (p *Polipo[T]) Do(ctx context.Context) ([]T, error) {
	if len(p.tasks) == 0 {
		return nil, errors.New("no tasks to do")
	}

	resultsChan := make(chan result[T])
	wg := sync.WaitGroup{}

	wg.Add(len(p.tasks))

	for _, task := range p.tasks {
		go func(t Task[T]) {
			defer wg.Done()
			items, err := t()

			select {
			case resultsChan <- result[T]{items, err}:
			case <-ctx.Done():
			}
		}(task)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var (
		results []T
		errs    []error
	)

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
