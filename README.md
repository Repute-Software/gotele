<<<<<<< HEAD
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
=======
# Gotele - Telegram Bot API Client for Go

A comprehensive, production-ready Telegram Bot API client library for Go with enterprise-grade features.

## Features

- ✅ **Enhanced Error Handling** - Comprehensive error types with retry logic support
- ✅ **Context Support** - Full context support for cancellation and timeouts
- ✅ **Comprehensive Type System** - Complete Telegram Bot API types
- ✅ **File Upload & Media** - Support for all media types with validation
- ✅ **Webhook Support** - Real-time message processing with security
- ✅ **Production Ready** - Extensive testing and error handling

## Quick Start
>>>>>>> af0c0f5 (Made some changes.)

```go
package main

import (
    "log"
    gotele "github.com/repute-software/gotele/telegram"
)

func main() {
<<<<<<< HEAD
    token := os.Getenv("BOT_DEV_TOKEN") // or load from config
    bot := gotele.NewBot(token)
    if err := bot.SendMessage(123456789, "Hello from gotele!"); err != nil {
=======
    // Create bot
    bot := gotele.NewBot("your_bot_token")
    
    // Send a message
    err := bot.SendMessage(chatID, "Hello, World!")
    if err != nil {
>>>>>>> af0c0f5 (Made some changes.)
        log.Fatal(err)
    }
}
```

<<<<<<< HEAD
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

=======
## Installation

```bash
go get github.com/repute-software/gotele/telegram
```

## Examples

See the [examples directory](./examples/) for comprehensive usage examples:

- [Context Support](./examples/context/) - Request cancellation and timeouts
- [Advanced Types](./examples/advanced_types/) - Interactive keyboards and formatting
- [File Upload](./examples/file_upload/) - Media handling and file operations
- [Webhook Support](./examples/webhook/) - Real-time message processing

## Documentation

- [API Reference](https://pkg.go.dev/github.com/repute-software/gotele/telegram)
- [Examples](./examples/README.md)
- [Telegram Bot API](https://core.telegram.org/bots/api)

## License

MIT License - see [LICENSE](./LICENSE) file for details.
>>>>>>> af0c0f5 (Made some changes.)
