package main

import (
	"context"
	"fmt"
	"time"

	"polipo"
)

type Fish struct {
	Name string
}

func main() {
	o := polipo.NewPolipo[Fish]()

	o.AddTentacle(func() ([]Fish, error) {
		return []Fish{
			{Name: "Salmon"},
			{Name: "Tuna"},
			{Name: "Trout"},
			{Name: "Cod"},
			{Name: "Haddock"},
			{Name: "Mackerel"},
		}, nil
	})

	o.AddTentacle(func() ([]Fish, error) {
		return nil, nil
	})

	o.AddTentacle(func() ([]Fish, error) {
		return []Fish{
			{Name: "Swordfish"},
			{Name: "Marlin"},
			{Name: "Barracuda"},
			{Name: "Mahi Mahi"},
			{Name: "Wahoo"},
			{Name: "Kingfish"},
		}, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results, err := o.Do(ctx)
	if err != nil {
		panic(err) // this is just an example, don't panic in production
	}

	fmt.Printf("Results: %+v\n", results)
}
