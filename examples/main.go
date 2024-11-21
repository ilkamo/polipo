package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ilkamo/polipo"
)

type TaskResult struct {
	FishName string
}

func main() {
	p := polipo.NewPolipo[TaskResult]()

	p.AddTask(func() ([]TaskResult, error) {
		return []TaskResult{
			{FishName: "Salmon"},
			{FishName: "Tuna"},
			{FishName: "Trout"},
			{FishName: "Cod"},
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results, err := p.Do(ctx)
	if err != nil {
		panic(err) // this is just an example, don't panic in production
	}

	fmt.Printf("Results: %+v\n", results)
}
