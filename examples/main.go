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

	p.AddTask(func() (TaskResult, error) {
		return TaskResult{
			Fishes: []string{
				"Salmon",
				"Tuna",
				"Trout",
				"Cod",
			},
		}, nil
	})

	p.AddTask(func() (TaskResult, error) {
		return TaskResult{}, nil
	})

	p.AddTask(func() (TaskResult, error) {
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
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results, err := p.Do(ctx)
	if err != nil {
		panic(err) // this is just an example, don't panic in production
	}

	fmt.Printf("Results: %+v\n", results)
}
