# go-puush-client

This repository contains a go implementation of the puush api. My plan for the future is to add a cross-platform desktop application, which will aim to recreate the original client as closely as possible. Right now, the repository only contains the barebones api implementation.

## Usage example

The `pkg/puush` package contains the implementation of the puush api.
Here's a simple example of how to upload a file:

```go
package main

import (
    "fmt"
    "github.com/Lekuruu/go-puush-client/pkg/puush"
)

func main() {
	// Login with email & password
	client := NewClientFromLogin("you@example.com", "example-password")
	client.SetBaseURL("https://puush.me")

	// Or login with an api key
    // client := NewClientFromApiKey("you@example.com", "api-key")

	if err := client.Authenticate(); err != nil {
		panic(err)
	}

	file, err := os.Open("example.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	url, err := client.Upload(file)
    if err != nil {
        panic(err)
    }

    fmt.Printf("File uploaded successfully: %s\n", url)
}
```
