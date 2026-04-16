package polipo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithMaxConcurrency(t *testing.T) {
	type result struct{}

	t.Run("should set the maximum number of concurrent tasks to run", func(t *testing.T) {
		maxConcurrency := 5
		p := NewPolipo[result](
			WithMaxConcurrency[result](maxConcurrency),
		)
		require.Equal(t, maxConcurrency, p.maxConcurrency)
	})

	t.Run("should clamp zero to 1", func(t *testing.T) {
		p := NewPolipo[result](
			WithMaxConcurrency[result](0),
		)
		require.Equal(t, 1, p.maxConcurrency)
	})

	t.Run("should clamp negative to 1", func(t *testing.T) {
		p := NewPolipo[result](
			WithMaxConcurrency[result](-5),
		)
		require.Equal(t, 1, p.maxConcurrency)
	})
}
