package polipo

import (
	"context"
	"errors"
	"sync"
)

// Tentacle is a function that can be run by the Polipo.
// Arguments could be passed in as a closure.
// The function should return a slice of generic items and an error.
type Tentacle[T any] func() ([]T, error)

// Polipo stores a list of Tentacles to be run concurrently.
type Polipo[T any] struct {
	wg        sync.WaitGroup
	tentacles []Tentacle[T]
}

// NewPolipo creates a new Polipo.
func NewPolipo[T any]() Polipo[T] {
	return Polipo[T]{
		wg:        sync.WaitGroup{},
		tentacles: make([]Tentacle[T], 0),
	}
}

// AddTentacle adds a Tentacle to the Polipo. The Tentacle function will be run when Do is called.
// Consider this as adding a function to a list of functions to be run.
func (p *Polipo[T]) AddTentacle(tentacle Tentacle[T]) {
	p.tentacles = append(p.tentacles, tentacle)
}

// Do executes all the Tentacles functions concurrently.
// This is a blocking function that will return when all the Tentacles have finished their work.
func (p *Polipo[T]) Do(ctx context.Context) ([]T, error) {
	if len(p.tentacles) == 0 {
		return nil, errors.New("no tentacles to catch")
	}

	resultsChan := make(chan catchResult[T])

	p.wg.Add(len(p.tentacles))

	for _, tentacle := range p.tentacles {
		go func(tentacle Tentacle[T]) {
			defer p.wg.Done()
			items, err := tentacle()

			select {
			case resultsChan <- catchResult[T]{items, err}:
			case <-ctx.Done():
			}
		}(tentacle)
	}

	go func() {
		p.wg.Wait()
		close(resultsChan)
	}()

	var (
		results      []T
		polipoErrors []error
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
				return results, errors.Join(polipoErrors...)
			}

			if res.err != nil {
				polipoErrors = append(polipoErrors, res.err)
			}

			results = append(results, res.items...)
		case <-ctx.Done():
			return results, errors.Join(append(polipoErrors, ctx.Err())...)
		}
	}
}

type catchResult[T any] struct {
	items []T
	err   error
}
