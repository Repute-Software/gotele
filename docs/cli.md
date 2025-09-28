## CLI (demo) in `cmd/gotele`

The repository includes a small demo CLI showcasing context timeouts and update retrieval.

### Run directly

```bash
export BOT_DEV_TOKEN=your_token
go run ./cmd/gotele
```

### Install

```bash
go install ./cmd/gotele
gotele
```

The CLI demonstrates:

- Loading token from environment
- Creating a bot with a custom timeout
- Getting updates with and without context

