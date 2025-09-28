## Usage

### Create a bot

```go
bot := gotele.NewBot(os.Getenv("BOT_DEV_TOKEN"))
```

With a custom timeout:

```go
bot := gotele.NewBotWithTimeout(token, 15*time.Second)
```

### Send a message

```go
_ = bot.SendMessage(chatID, "Hello!")
```

With options:

```go
opts := &gotele.SendMessageOptions{ParseMode: "Markdown"}
_ = bot.SendMessageAdvanced(chatID, "*bold* _italics_", opts)
```

With context:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
_ = bot.SendMessageAdvancedWithContext(ctx, chatID, "hello", nil)
```

### Edit and delete

```go
_ = bot.EditMessageText(&gotele.EditMessageTextOptions{ChatID: chatID, MessageID: 123, Text: "Updated"})
_ = bot.DeleteMessage(chatID, 123)
```

### Receive updates (long polling)

```go
updates, err := bot.GetUpdates(0)
_ = err
for _, u := range updates {
    if u.Message != nil {
        // handle message
    }
}
```

### Keyboards and entities

```go
keyboard := gotele.InlineKeyboardMarkup{InlineKeyboard: [][]gotele.InlineKeyboardButton{{{Text: "Click", CallbackData: "cb"}}}}
_ = bot.SendMessageAdvanced(chatID, "Choose:", &gotele.SendMessageOptions{ReplyMarkup: keyboard})
```

