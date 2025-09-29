# Gotele Examples

This directory contains example programs demonstrating various features of the gotele Telegram API client library.

## Prerequisites

Before running any examples, you need to:

1. Set up a Telegram bot by talking to [@BotFather](https://t.me/botfather)
2. Get your bot token
3. Set the `BOT_DEV_TOKEN` environment variable

```bash
export BOT_DEV_TOKEN="your_bot_token_here"
```

## Available Examples

### 1. Context Examples (`context/`)
Demonstrates context support for request cancellation and timeouts.

**Features:**
- Context with timeout
- Context cancellation
- Custom timeout configuration
- Context with values

**Run:**
```bash
go run ./examples/context
```

### 2. Advanced Types Examples (`advanced_types/`)
Shows how to use advanced message types and interactive elements.

**Features:**
- Inline keyboards
- Reply keyboards
- Message entities and formatting
- Photo sending with captions
- Message editing

**Run:**
```bash
go run ./examples/advanced_types
```

### 3. File Upload Examples (`file_upload/`)
Demonstrates file upload and media handling capabilities.

**Features:**
- Document uploads
- Video uploads with thumbnails
- Audio file uploads
- Media groups (albums)
- File download
- File validation

**Run:**
```bash
go run ./examples/file_upload
```

### 4. Webhook Examples (`webhook/`)
Shows webhook setup and real-time message processing.

**Features:**
- Webhook setup and management
- HTTP/HTTPS webhook servers
- Update routing and handling
- Signature validation
- Middleware support

**Run:**
```bash
go run ./examples/webhook
```

### 5. Edit Messages & Keyboards (`edit_and_keyboard/`)
Demonstrates message editing and reply keyboard functionality.

**Features:**
- Edit message text
- Edit message captions  
- Edit reply markup
- Reply keyboard creation
- Advanced keyboard buttons
- One-time keyboards

**Run:**
```bash
go run ./examples/edit_and_keyboard
```

## Configuration

Most examples require a valid chat ID to send messages to. You can:

1. Start a conversation with your bot
2. Send a message to get your chat ID
3. Update the examples with your chat ID

For webhook examples, you'll need:
- A publicly accessible HTTPS URL
- A valid SSL certificate
- Port forwarding if running locally

## Notes

- Examples are designed for demonstration purposes
- Some examples may require additional setup (like webhook URLs)
- Replace placeholder values (like chat IDs) with actual values
- Examples include error handling but may not cover all edge cases

## Troubleshooting

**Common issues:**

1. **"BOT_DEV_TOKEN is not set"**
   - Make sure you've set the environment variable
   - Verify your bot token is correct

2. **"chat not found" errors**
   - Make sure you've started a conversation with your bot
   - Use the correct chat ID

3. **Webhook errors**
   - Ensure your webhook URL is publicly accessible
   - Check that you're using HTTPS (required by Telegram)
   - Verify your SSL certificate is valid

4. **File upload errors**
   - Check file paths and permissions
   - Ensure files don't exceed Telegram's size limits
   - Verify file formats are supported
