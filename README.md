# Dwaz - Discord Wrapper API for Zwafriya
> Dwaz 5alina men lmar9a kat3reb tsaybe bot b dwaz?

## Overview

Dwaz is a lightweight and modern Discord API wrapper written in Go, designed for developers building Discord bots or integrations. With a focus on simplicity and performance, Dwaz provides an intuitive interface to interact with Discord's API, leveraging Go's concurrency features for efficient bot development. The name "Zwafriya," inspired by Moroccan Arabic.

## Installation

To use Dwaz in your Go project, install it via:

```bash
go get github.com/marouanesouiri/dwaz
```

Ensure you have Go 1.22 or later installed, as specified in the project’s `go.mod`.

## Usage

Here’s a basic Ping Pong example to get started with Dwaz:

```go
package main

import (
    "github.com/marouanesouiri/dwaz"
    "context"
    "fmt"
)

func main() {
    // Initialize a new Dwaz client
    client, err := dwaz.New(
		dwaz.WithToken("YOUR_BOT_TOKEN"),
		dwaz.WithIntents(dwaz.GatewayIntentGuildMessages, dwaz.GatewayIntentMessageContent),
    )

    // Add message create even handlers
    client.OnMessageCreate(func(event *dwaz.MessageCreateEvent) {
        if event.Message.Content == "!ping" {
            fmt.Println("Pong!")
        }
    })

    // Start the bot
    client.Start(context.TODO())
}
```

Replace `YOUR_BOT_TOKEN` with your Discord bot token. Check the [documentation](https://pkg.go.dev/github.com/marouanesouiri/dwaz) for more examples and API details.

## Badges

[![Go Reference](https://pkg.go.dev/badge/github.com/marouanesouiri/dwaz.svg)](https://pkg.go.dev/github.com/marouanesouiri/dwaz)
[![Go Report](https://goreportcard.com/badge/github.com/marouanesouiri/dwaz)](https://goreportcard.com/report/github.com/marouanesouiri/dwaz)
[![Go Version](https://img.shields.io/github/go-mod/go-version/marouanesouiri/dwaz)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/marouanesouiri/dwaz/blob/master/LICENSE)
[![Yada Version](https://img.shields.io/github/v/tag/marouanesouiri/dwaz?label=release)](https://github.com/marouanesouiri/dwaz/releases/latest)
[![Issues](https://img.shields.io/github/issues/marouanesouiri/dwaz)](https://github.com/marouanesouiri/dwaz/issues)
[![Last Commit](https://img.shields.io/github/last-commit/marouanesouiri/dwaz)](https://github.com/marouanesouiri/dwaz/commits/main)
[![Lines of Code](https://tokei.rs/b1/github/marouanesouiri/dwaz)](https://github.com/marouanesouiri/dwaz)
