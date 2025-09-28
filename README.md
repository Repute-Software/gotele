# gotele

Lightweight, pragmatic Go client for the Telegram Bot API with first-class context support, simple file uploads, and ergonomic helpers for webhooks and interactive features.

## Features

- Context-aware methods for all requests (`WithContext` variants)
- Simple message APIs: send, edit, delete
- Rich options: parse modes, entities, reply markup (inline/reply keyboards)
- File uploads (documents, audio, video) and media groups (albums)
- File download helpers and size validation
- Webhook setup, HTTP/TLS servers, signature validation, and middleware
- Backward-compatible long polling via `GetUpdates`

## Installation

```bash
go get github.com/repute-software/gotele/telegram
```

Import in your code:

```go
import gotele "github.com/repute-software/gotele/telegram"
```

## Quickstart

```go
package main

import (
    "log"
    gotele "github.com/repute-software/gotele/telegram"
)

func main() {
    token := os.Getenv("BOT_DEV_TOKEN") // or load from config
    bot := gotele.NewBot(token)
    if err := bot.SendMessage(123456789, "Hello from gotele!"); err != nil {
        log.Fatal(err)
    }
}
```

More examples are available under `examples/` and in the docs below.

## Documentation

- Getting started and installation: `docs/installation.md`
- Usage guide (messages, options, updates): `docs/usage.md`
- Webhooks (servers, middleware, routing): `docs/webhooks.md`
- File uploads and downloads: `docs/file-uploads.md`
- Advanced types and interactive features: `docs/advanced-types.md`
- CLI helper (demo app in `cmd/gotele`): `docs/cli.md`
- Examples walkthrough: `docs/examples.md`
- Troubleshooting: `docs/troubleshooting.md`

## CLI (optional)

A small demo CLI lives in `cmd/gotele/`. You can run it with:

```bash
export BOT_DEV_TOKEN=your_token
go run ./cmd/gotele
```

Or install the binary:

```bash
go install ./cmd/gotele
gotele
```

## Requirements

- Go 1.20+
- A Telegram bot token from `@BotFather`

## Examples

See `examples/` for runnable programs showing contexts, advanced message types, file uploads, and webhooks. Start with `examples/README.md`.

## License

This project is licensed under the terms of the `LICENSE` file.

