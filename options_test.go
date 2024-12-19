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
}
