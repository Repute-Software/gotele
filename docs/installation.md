## Installation

### Prerequisites

- Go 1.20+
- Telegram bot token from `@BotFather`

### Install the library

```bash
go get github.com/repute-software/gotele/telegram
```

### Import

```go
import gotele "github.com/repute-software/gotele/telegram"
```

### Verify

```bash
go env GOPATH
```

Ensure your project builds:

```bash
go build ./...
```

### Configure token

Set `BOT_DEV_TOKEN` for local development:

```bash
export BOT_DEV_TOKEN="your_bot_token"
```

