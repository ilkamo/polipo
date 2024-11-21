package polipo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"polipo"
)

type Fish struct {
	Name string
}

func TestPolipo_Catch(t *testing.T) {
	t.Run("should return a list of all fishes", func(t *testing.T) {
		ctx := context.TODO()

		expectedFishes := []Fish{
			{Name: "Salmon"},
			{Name: "Tuna"},
			{Name: "Trout"},
			{Name: "Cod"},
			{Name: "Haddock"},
			{Name: "Mackerel"},
			{Name: "Swordfish"},
			{Name: "Marlin"},
			{Name: "Barracuda"},
			{Name: "Mahi Mahi"},
			{Name: "Wahoo"},
			{Name: "Kingfish"},
		}

		g := polipo.NewPolipo[Fish]()

		g.AddTentacle(func() ([]Fish, error) {
			return []Fish{
				{Name: "Salmon"},
				{Name: "Tuna"},
				{Name: "Trout"},
				{Name: "Cod"},
				{Name: "Haddock"},
				{Name: "Mackerel"},
			}, nil
		})

		g.AddTentacle(func() ([]Fish, error) {
			return nil, nil
		})

		g.AddTentacle(func() ([]Fish, error) {
			return []Fish{
				{Name: "Swordfish"},
				{Name: "Marlin"},
				{Name: "Barracuda"},
				{Name: "Mahi Mahi"},
				{Name: "Wahoo"},
				{Name: "Kingfish"},
			}, nil
		})

		fishes, err := g.Do(ctx)
		require.NoError(t, err)
		require.Len(t, fishes, 12)
		require.ElementsMatch(t, expectedFishes, fishes)
	})

	t.Run("should return an error if one of the tentacles returns an error", func(t *testing.T) {
		ctx := context.TODO()

		g := polipo.NewPolipo[Fish]()

		g.AddTentacle(func() ([]Fish, error) {
			return []Fish{
				{Name: "Swordfish"},
				{Name: "Marlin"},
			}, nil
		})

		g.AddTentacle(func() ([]Fish, error) {
			return nil, errors.New("nothing in the ocean")
		})

		fishes, err := g.Do(ctx)
		require.ErrorContains(t, err, "nothing in the ocean")
		require.Len(t, fishes, 2)
		require.ElementsMatch(t, []Fish{
			{Name: "Swordfish"},
			{Name: "Marlin"},
		}, fishes)
	})

	t.Run("should return an error if the context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		g := polipo.NewPolipo[Fish]()

		g.AddTentacle(func() ([]Fish, error) {
			return []Fish{
				{Name: "Swordfish"},
				{Name: "Marlin"},
			}, nil
		})

		cancel()

		fishes, err := g.Do(ctx)
		require.ErrorContains(t, err, "context canceled")
		require.Empty(t, fishes)
	})

	t.Run("should return an error if the context is canceled because of timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		g := polipo.NewPolipo[Fish]()

		g.AddTentacle(func() ([]Fish, error) {
			return []Fish{
				{Name: "Swordfish"},
			}, nil
		})

		g.AddTentacle(func() ([]Fish, error) {
			time.Sleep(time.Second * 10)
			return []Fish{
				{Name: "Marlin"},
			}, nil
		})

		fishes, err := g.Do(ctx)
		require.ErrorContains(t, err, "context deadline exceeded")
		require.ElementsMatch(t, []Fish{
			{Name: "Swordfish"},
		}, fishes)
	})
}

func BenchmarkPolipo_Catch(b *testing.B) {
	ctx := context.TODO()

	g := polipo.NewPolipo[Fish]()

	g.AddTentacle(func() ([]Fish, error) {
		return []Fish{
			{Name: "Salmon"},
			{Name: "Tuna"},
			{Name: "Trout"},
			{Name: "Cod"},
			{Name: "Haddock"},
			{Name: "Mackerel"},
		}, nil
	})

	g.AddTentacle(func() ([]Fish, error) {
		return nil, nil
	})

	g.AddTentacle(func() ([]Fish, error) {
		return []Fish{
			{Name: "Swordfish"},
			{Name: "Marlin"},
			{Name: "Barracuda"},
			{Name: "Mahi Mahi"},
			{Name: "Wahoo"},
			{Name: "Kingfish"},
		}, nil
	})

	for i := 0; i < b.N; i++ {
		_, _ = g.Do(ctx)
	}
}
