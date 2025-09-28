## Troubleshooting

### BOT_DEV_TOKEN is not set

- Export `BOT_DEV_TOKEN` and verify the value
- Confirm the token with `@BotFather`

### chat not found / 400 errors

- Start a conversation with your bot first
- Use the correct numeric `chat_id`

### Webhook errors

- Ensure your URL is publicly reachable over HTTPS
- Provide a valid certificate or use TLS server variant
- Check `SecretToken` header and validation

### File upload errors

- Verify file paths and permissions
- Check size limits with `ValidateFileSize`
- Ensure supported MIME types/extensions

### Timeouts

- Increase bot timeout via `NewBotWithTimeout`
- Use `WithContext` methods and extend deadlines

