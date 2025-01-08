# Reitti

Reitti[^1] is a simple and flexible router package for Go. It allows you to define routes with template parameters and match them against paths. It can be used outside HTTP APIs.

[^1]: "reitti" is Finnish for "route".

## Installation

To install Reitti, use `go get`:

```sh
go get github.com/wolverian/reitti
```

## Usage

Here is a simple example of how to use Reitti:

```go
package main

import (
	"context"
	"fmt"

	"github.com/wolverian/reitti"
)

func main() {
	r := &reitti.Router{}
	r.Add("repos/{owner}/{repo}/issues", func(ctx context.Context, owner, repo string) (any, error) {
		return fmt.Sprintf("Owner: %s, Repo: %s", owner, repo), nil
	})

	handler, err := r.Match("repos/wolverian/reitti/issues")
	if err != nil {
		fmt.Println("No route found")
		return
	}

	result, err := handler(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(result) // Output: Owner: wolverian, Repo: reitti
}
```
