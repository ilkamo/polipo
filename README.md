# polipo

## Overview

polipo is a Go library designed to manage and execute concurrent tasks using a flexible and extensible approach. The
library allows you to define and add multiple tasks (tentacles) to a `Polipo` instance and execute them concurrently,
handling results and errors efficiently.

It is particularly useful when you need to execute multiple tasks concurrently and collect the results. For example, you
can use it to fetch data from multiple sources or providers, process it concurrently, and aggregate the results into a single output.

The name "polipo" is derived from the Italian word for "octopus," which has multiple tentacles that can perform tasks independently.

## Features

- **Generic Support**: The library uses Go generics to support any data type.
- **Concurrent Execution**: Tasks are executed concurrently, leveraging Go's goroutines.
- **Context Support**: Execution can be controlled and canceled using Go's `context.Context`.
- **Error Handling**: Collects and returns errors from all tasks.

## Installation

To install the library, use `go get`:

```sh
go get github.com/ilkamo/polipo
```

## Usage

### Creating a Polipo instance

To create a new `Polipo` instance, specify the type of data it will handle:

```go
import "github.com/ilkamo/polipo"

type TaskResult struct {
ID   int
Name string
}

p := polipo.NewPolipo[TaskResult]()
```

### Adding Tentacles

Add tentacles (tasks) to the `Polipo` instance. Each tentacle is a function that returns a slice of items and an error:

```go
p.AddTentacle(func () ([]TaskResult, error) {
return []TaskResult{
{ID: 1, Name: "Task1"},
{ID: 2, Name: "Task2"},
}, nil
})
```

### Run Tentacles

Run all tentacles concurrently using the `Do` method. Pass a `context.Context` to control execution:

```go
ctx := context.TODO()

results, err := p.Do(ctx)
if err != nil {
log.Fatal(err)
}

for _, result := range results {
fmt.Println(result.Name)
}
```

## Example

Here is a complete example:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ilkamo/polipo"
)

type TaskResult struct {
	ID   int
	Name string
}

func main() {
	ctx := context.TODO()
	p := polipo.NewPolipo[TaskResult]()

	p.AddTentacle(func() ([]TaskResult, error) {
		return []TaskResult{
			{ID: 1, Name: "Task1"},
			{ID: 2, Name: "Task2"},
		}, nil
	})

	p.AddTentacle(func() ([]TaskResult, error) {
		return []TaskResult{
			{ID: 3, Name: "Task3"},
			{ID: 4, Name: "Task4"},
		}, nil
	})

	results, err := p.Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result.Name)
	}
}
```

## Testing

To run tests, use the following command:

```sh
make test
```

## Benchmarking

To run benchmarks, use the following command:

```sh
make benchmark
```

## Linting

To lint the code, use the following command:

```sh
make lint-fix
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
