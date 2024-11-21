package polipo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/polipo"
)

type TaskResult struct {
	FishName string
}

func TestPolipo_Do(t *testing.T) {
	t.Run("should return a list of all fishes", func(t *testing.T) {
		ctx := context.TODO()

		expectedFishes := []TaskResult{
			{FishName: "Salmon"},
			{FishName: "Tuna"},
			{FishName: "Trout"},
			{FishName: "Cod"},
			{FishName: "Haddock"},
			{FishName: "Mackerel"},
			{FishName: "Swordfish"},
			{FishName: "Marlin"},
			{FishName: "Barracuda"},
			{FishName: "Mahi Mahi"},
			{FishName: "Wahoo"},
			{FishName: "Kingfish"},
		}

		p := polipo.NewPolipo[TaskResult]()

		p.AddTask(func() ([]TaskResult, error) {
			return []TaskResult{
				{FishName: "Salmon"},
				{FishName: "Tuna"},
				{FishName: "Trout"},
				{FishName: "Cod"},
				{FishName: "Haddock"},
				{FishName: "Mackerel"},
			}, nil
		})

		p.AddTask(func() ([]TaskResult, error) {
			return nil, nil
		})

		p.AddTask(func() ([]TaskResult, error) {
			return []TaskResult{
				{FishName: "Swordfish"},
				{FishName: "Marlin"},
				{FishName: "Barracuda"},
				{FishName: "Mahi Mahi"},
				{FishName: "Wahoo"},
				{FishName: "Kingfish"},
			}, nil
		})

		fishes, err := p.Do(ctx)
		require.NoError(t, err)
		require.Len(t, fishes, 12)
		require.ElementsMatch(t, expectedFishes, fishes)
	})

	t.Run("should return an error if one of the tasks returns an error", func(t *testing.T) {
		ctx := context.TODO()

		p := polipo.NewPolipo[TaskResult]()

		p.AddTask(func() ([]TaskResult, error) {
			return []TaskResult{
				{FishName: "Swordfish"},
				{FishName: "Marlin"},
			}, nil
		})

		p.AddTask(func() ([]TaskResult, error) {
			return nil, errors.New("nothing in the ocean")
		})

		fishes, err := p.Do(ctx)
		require.ErrorContains(t, err, "nothing in the ocean")
		require.Len(t, fishes, 2)
		require.ElementsMatch(t, []TaskResult{
			{FishName: "Swordfish"},
			{FishName: "Marlin"},
		}, fishes)
	})

	t.Run("should return an error if the context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		p := polipo.NewPolipo[TaskResult]()

		p.AddTask(func() ([]TaskResult, error) {
			return []TaskResult{
				{FishName: "Swordfish"},
				{FishName: "Marlin"},
			}, nil
		})

		cancel()

		fishes, err := p.Do(ctx)
		require.ErrorContains(t, err, "context canceled")
		require.Empty(t, fishes)
	})

	t.Run("should return an error if the context is canceled because of timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		p := polipo.NewPolipo[TaskResult]()

		p.AddTask(func() ([]TaskResult, error) {
			return []TaskResult{
				{FishName: "Swordfish"},
			}, nil
		})

		p.AddTask(func() ([]TaskResult, error) {
			time.Sleep(time.Second * 10)
			return []TaskResult{
				{FishName: "Marlin"},
			}, nil
		})

		fishes, err := p.Do(ctx)
		require.ErrorContains(t, err, "context deadline exceeded")
		require.ElementsMatch(t, []TaskResult{
			{FishName: "Swordfish"},
		}, fishes)
	})
}

func BenchmarkPolipo_Do(b *testing.B) {
	ctx := context.TODO()

	p := polipo.NewPolipo[TaskResult]()

	p.AddTask(func() ([]TaskResult, error) {
		return []TaskResult{
			{FishName: "Salmon"},
			{FishName: "Tuna"},
			{FishName: "Trout"},
			{FishName: "Cod"},
			{FishName: "Haddock"},
			{FishName: "Mackerel"},
		}, nil
	})

	p.AddTask(func() ([]TaskResult, error) {
		return nil, nil
	})

	p.AddTask(func() ([]TaskResult, error) {
		return []TaskResult{
			{FishName: "Swordfish"},
			{FishName: "Marlin"},
			{FishName: "Barracuda"},
			{FishName: "Mahi Mahi"},
			{FishName: "Wahoo"},
			{FishName: "Kingfish"},
		}, nil
	})

	for i := 0; i < b.N; i++ {
		_, _ = p.Do(ctx)
	}
}
