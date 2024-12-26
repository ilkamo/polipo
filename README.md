# Polipo

### Go library designed to manage and execute concurrent tasks using generics
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://pkg.go.dev/badge/github.com/ilkamo/polipo?status.svg)](https://pkg.go.dev/github.com/ilkamo/polipo?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/ilkamo/polipo)](https://goreportcard.com/report/ilkamo/polipo)

**Polipo** library allows you to define multiple tasks and execute them concurrently, handling results and errors efficiently. For example, you can use it to fetch data from multiple sources or providers or to perform multiple calculations in parallel and combine the results into a single output.

Unlike other similar libraries, **polipo** uses channels under the hood. No mutexes or locks are used, which makes it faster and more efficient.

<img src="assets/polipo.webp" width="400">

The name "_polipo_" is derived from the Italian word for "_octopus_", which has multiple tentacles that can perform tasks independently. This is analogous to the library's ability to execute multiple tasks concurrently.

## Features

- **Generic Support**: The library uses Go generics to support any data type.
- **Concurrent Execution**: Tasks are executed concurrently, leveraging Go's goroutines. The number of max concurrent tasks
  can be controlled by using the `WithMaxConcurrency` option.
- **Context Support**: Execution can be controlled and canceled using Go's `context.Context`.
- **Error Handling**: Collects and returns errors from all tasks.

## Installation

To install the library, use `go get`:

```sh
go get github.com/ilkamo/polipo
```

## Usage

### Creating a Polipo instance

To create a new `Polipo` instance, specify the type of data it will handle. It can be any data type, such as a struct or a primitive type:

```go
import "github.com/ilkamo/polipo"

type TaskResult struct {
    ID   int
    Name string
}

p := polipo.NewPolipo[TaskResult]()
```

### Adding Tasks

Each task is a function that returns a slice of items and an error:

```go
p := polipo.NewPolipo[TaskResult]()

err := p.AddTask(func () (TaskResult, error) {
    return TaskResult{ID: 1, Name: "Task1"}, nil
})
```

### Run Tasks

Run all tasks concurrently using the `Do` method. Pass a `context.Context` to control the execution:

```go
ctx := context.TODO()

p := polipo.NewPolipo[TaskResult]()

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

    p.AddTask(func() (TaskResult, error) {
      return TaskResult{ID: 1, Name: "Task1"}, nil
    })

    p.AddTask(func() (TaskResult, error) {
      return TaskResult{ID: 2, Name: "Task2"}, nil
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
