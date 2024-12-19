package polipo

type Option[T any] func(*Polipo[T])

// WithMaxConcurrency sets the maximum number of concurrent tasks to run.
func WithMaxConcurrency[T any](maxConcurrency int) Option[T] {
	return func(p *Polipo[T]) {
		p.maxConcurrency = maxConcurrency
	}
}
