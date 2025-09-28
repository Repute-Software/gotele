## Advanced types and interactive features

### Inline keyboards

```go
keyboard := gotele.InlineKeyboardMarkup{
    InlineKeyboard: [][]gotele.InlineKeyboardButton{
        {{Text: "Open", URL: "https://example.com"}},
        {{Text: "Callback", CallbackData: "cb:data"}},
    },
}
_ = bot.SendMessageAdvanced(chatID, "Choose:", &gotele.SendMessageOptions{ReplyMarkup: keyboard})
```

### Reply keyboards

```go
rk := gotele.ReplyKeyboardMarkup{Keyboard: [][]gotele.KeyboardButton{{{Text: "Yes"}, {Text: "No"}}}}
_ = bot.SendMessageAdvanced(chatID, "Reply:", &gotele.SendMessageOptions{ReplyMarkup: rk})
```

### Entities and formatting

```go
entities := []gotele.MessageEntity{{Type: "bold", Offset: 0, Length: 4}}
_ = bot.SendMessageAdvanced(chatID, "Bold text", &gotele.SendMessageOptions{Entities: entities})
```

### Edit message

```go
_ = bot.EditMessageText(&gotele.EditMessageTextOptions{ChatID: chatID, MessageID: msgID, Text: "Updated"})
```

### Answer callback queries and inline queries

```go
_ = bot.AnswerCallbackQuery(&gotele.AnswerCallbackQueryOptions{CallbackQueryID: id, Text: "Done"})

results := []gotele.InlineQueryResult{{Type: "article", ID: "1", Title: "Hi", InputMessageContent: gotele.InputTextMessageContent{MessageText: "Hello"}}}
_ = bot.AnswerInlineQuery(&gotele.AnswerInlineQueryOptions{InlineQueryID: qid, Results: results})
```

