package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ilkamo/polipo"
)

type TaskResult struct {
	Fishes []string
}

func main() {
	p := polipo.NewPolipo[TaskResult]()

	if err := p.AddTask(func() (TaskResult, error) {
		return TaskResult{
			Fishes: []string{
				"Salmon",
				"Tuna",
				"Trout",
				"Cod",
			},
		}, nil
	}); err != nil {
		panic(err) // this is just an example, don't panic in production
	}

	if err := p.AddTask(func() (TaskResult, error) {
		return TaskResult{}, nil
	}); err != nil {
		panic(err)
	}

	if err := p.AddTask(func() (TaskResult, error) {
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
	}); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results, err := p.Do(ctx)
	if err != nil {
		panic(err) // this is just an example, don't panic in production
	}

	fmt.Printf("Results: %+v\n", results)
}
