## Webhooks

### Set a webhook

```go
options := &gotele.SetWebhookOptions{
    URL:         "https://your-domain.com/webhook",
    SecretToken: "your-secret",
}
_ = bot.SetWebhook(options)
```

### Start a server

```go
server := &gotele.WebhookServer{
    Port:        "8080",
    SecretToken: "your-secret",
    Handler: func(u *gotele.Update) error {
        // route and handle
        return nil
    },
}
_ = bot.StartWebhookServer(server)
```

TLS variant:

```go
_ = bot.StartWebhookServerTLS(server, "cert.pem", "key.pem")
```

### Middleware and utilities

- `WebhookMiddleware(secret, next)` to validate signatures and pass through
- `WebhookLogger(next)` to log requests
- `ValidateWebhookSignature(secret, body, signature)` for manual checks
- `ProcessWebhookUpdate(update, handlers)` for typed routing

