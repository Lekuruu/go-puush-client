# go-puush-client

This repository contains a go implementation of the puush api and an almost feature-complete, cross-platform desktop application recreation of the original puush client.
My main motivation was to have a usable puush client for linux, since it was the only platform that didn't have an official client. It was also a good way for me to learn more about [fyne](https://fyne.io/), which is a cross-platform gui toolkit for go.

<img src="https://raw.githubusercontent.com/Lekuruu/go-puush-client/refs/heads/main/.github/screenshot-1.png" alt="tray" width="400px" />
<img src="https://raw.githubusercontent.com/Lekuruu/go-puush-client/refs/heads/main/.github/screenshot-2.png" alt="notification" width="400px" />

## Build instructions

If you know at least a little bit of how to use go, this should be pretty straightforward.

```bash
go build -o puush-client ./cmd/desktop
```

Since this project uses [fyne](https://fyne.io/), you may also need some C compiler dependencies depending on your platform. See [Fyne's prerequisite documentation](https://docs.fyne.io/started/quick/) for more information.  
Note that these can be pretty painful to set up on windows especially (speaking from experience).

## Screenshot Providers

For macOS and Windows, the application uses native implementations for taking screenshots, so you won't have to worry about installing any dependencies.  
On Linux, however, you may have to install one depending on your distribution:

- `spectacle` (Highly recommended for KDE)
- `gnome-screenshot` (Highly recommended for GNOME)

These two are the most feature complete and work flawlessly from testing.  
Other providers include `flameshot`, `maim` & `grim`/`slurp`, however, your mileage may vary with these.

## Progress

The main application is mostly feature complete, with a few minor exceptions:

- Context menu's + IPC
- An updater

Otherwise, the application does what I want it to do, which is to upload files and take screenshots.

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
