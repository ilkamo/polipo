package polipo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/polipo"
)

type TaskResult struct {
	Fishes []string
}

func TestPolipo_Do(t *testing.T) {
	t.Run("should return all fishes", func(t *testing.T) {
		tasks := []polipo.Task[TaskResult]{
			func() (TaskResult, error) {
				return TaskResult{
					Fishes: []string{
						"Salmon",
						"Tuna",
						"Trout",
						"Cod",
						"Haddock",
						"Mackerel",
					},
				}, nil
			},
			func() (TaskResult, error) {
				return TaskResult{}, nil
			},
			func() (TaskResult, error) {
				return TaskResult{
					Fishes: []string{
						"Swordfish",
						"Marlin",
						"Barracuda",
						"Mahi Mahi",
						"Wahoo",
						"Kingfish",
					},
				}, nil
			},
		}

		expected := []TaskResult{
			{
				Fishes: []string{
					"Salmon",
					"Tuna",
					"Trout",
					"Cod",
					"Haddock",
					"Mackerel",
				},
			},
			{
				Fishes: nil,
			},
			{
				Fishes: []string{
					"Swordfish",
					"Marlin",
					"Barracuda",
					"Mahi Mahi",
					"Wahoo",
					"Kingfish",
				},
			},
		}

		testCase := []struct {
			name        string
			concurrency int
		}{
			{
				name:        "max concurrency is 1",
				concurrency: 1,
			},
			{
				name:        "max concurrency is 5",
				concurrency: 5,
			},
			{
				name:        "max concurrency is 10",
				concurrency: 10,
			},
		}

		for _, tc := range testCase {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.TODO()

				p := polipo.NewPolipo[TaskResult]()

				for _, task := range tasks {
					err := p.AddTask(task)
					require.NoError(t, err)
				}

				allResults, err := p.Do(ctx)
				require.NoError(t, err)
				require.ElementsMatch(t, expected, allResults)
			})
		}
	})

	t.Run("should return an error if one of the tasks returns an error", func(t *testing.T) {
		ctx := context.TODO()

		p := polipo.NewPolipo[TaskResult]()

		err := p.AddTask(func() (TaskResult, error) {
			return TaskResult{
				Fishes: []string{
					"Swordfish",
					"Marlin",
				},
			}, nil
		})
		require.NoError(t, err)

		err = p.AddTask(func() (TaskResult, error) {
			return TaskResult{}, errors.New("nothing in the ocean")
		})
		require.NoError(t, err)

		allResults, err := p.Do(ctx)
		require.ErrorContains(t, err, "nothing in the ocean")
		require.Len(t, allResults, 1)
		require.ElementsMatch(t, []TaskResult{
			{
				Fishes: []string{
					"Swordfish",
					"Marlin",
				},
			},
		}, allResults)
	})

	t.Run("should return an error if the context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		p := polipo.NewPolipo[TaskResult]()

		err := p.AddTask(func() (TaskResult, error) {
			return TaskResult{
				Fishes: []string{
					"Swordfish",
					"Marlin",
				},
			}, nil
		})
		require.NoError(t, err)

		cancel()

		allResults, err := p.Do(ctx)
		require.ErrorContains(t, err, "context canceled")
		require.Empty(t, allResults)
	})

	t.Run("should return an error if the context is canceled because of timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		p := polipo.NewPolipo[TaskResult]()

		err := p.AddTask(func() (TaskResult, error) {
			return TaskResult{
				Fishes: []string{"Swordfish"},
			}, nil
		})
		require.NoError(t, err)

		err = p.AddTask(func() (TaskResult, error) {
			time.Sleep(time.Second * 10)
			return TaskResult{
				Fishes: []string{"Marlin"},
			}, nil
		})
		require.NoError(t, err)

		allResults, err := p.Do(ctx)
		require.ErrorContains(t, err, "context deadline exceeded")
		require.ElementsMatch(t, []TaskResult{
			{
				Fishes: []string{"Swordfish"},
			},
		}, allResults)
	})
}

var testCases = []struct {
	numberOfTasks int
}{
	{numberOfTasks: 100},
	{numberOfTasks: 1000},
	{numberOfTasks: 10000},
	{numberOfTasks: 100000},
}

func BenchmarkPolipo_Do(b *testing.B) {
	for _, tc := range testCases {
		b.Run(fmt.Sprintf("%d tasks", tc.numberOfTasks), func(b *testing.B) {
			ctx := context.TODO()

			p := polipo.NewPolipo[TaskResult]()

			for i := 0; i < tc.numberOfTasks; i++ {
				err := p.AddTask(func() (TaskResult, error) {
					return TaskResult{
						Fishes: []string{
							"Salmon",
							"Tuna",
							"Trout",
						},
					}, nil
				})
				require.NoError(b, err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = p.Do(ctx)
			}
		})
	}
}
