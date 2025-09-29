package main

import (
	"fmt"
	"log"
	"os"

	gotele "github.com/repute-software/gotele/telegram"
)

func main() {
	// Get bot token from environment
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	// Create bot
	bot := gotele.NewBot(botToken)
	chatID := int64(123456789) // Replace with your chat ID

	// Example 1: Send a message with a reply keyboard
	fmt.Println("=== Example 1: Reply Keyboard ===")

	// Create a simple keyboard
	keyboard := gotele.NewReplyKeyboard(
		gotele.NewKeyboardGrid(
			gotele.NewKeyboardRow(
				gotele.NewKeyboardButton("Option 1"),
				gotele.NewKeyboardButton("Option 2"),
			),
			gotele.NewKeyboardRow(
				gotele.NewKeyboardButton("Help"),
				gotele.NewKeyboardButton("Cancel"),
			),
		),
	)

	// Send message with keyboard
	messageOptions := &gotele.SendMessageOptions{
		ReplyMarkup: keyboard,
		ParseMode:   "HTML",
	}

	err := bot.SendMessageAdvanced(chatID, "Choose an option:", messageOptions)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		fmt.Println("Message sent successfully")
	}

	// Example 2: Edit message text
	fmt.Println("\n=== Example 2: Edit Message Text ===")

	// First send a message to edit
	err = bot.SendMessage(chatID, "This message will be edited")
	if err != nil {
		log.Printf("Failed to send initial message: %v", err)
		return
	}

	// For this example, we'll use a hardcoded message ID
	// In a real application, you'd get this from the message response
	messageID := 1 // Replace with actual message ID

	editOptions := &gotele.EditMessageTextOptions{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      "Message has been updated!",
		ParseMode: "HTML",
	}

	err = bot.EditMessageText(editOptions)
	if err != nil {
		log.Printf("Failed to edit message: %v", err)
	} else {
		fmt.Println("Message edited successfully")
	}

	// Example 3: Edit message caption (for media messages)
	fmt.Println("\n=== Example 3: Edit Message Caption ===")

	captionOptions := &gotele.EditMessageCaptionOptions{
		ChatID:    chatID,
		MessageID: messageID,
		Caption:   "This is an updated caption",
		ParseMode: "HTML",
	}

	err = bot.EditMessageCaption(captionOptions)
	if err != nil {
		log.Printf("Failed to edit caption: %v", err)
	} else {
		fmt.Println("Caption edited successfully")
	}

	// Example 4: Edit reply markup (change keyboard)
	fmt.Println("\n=== Example 4: Edit Reply Markup ===")

	// Create a new keyboard
	newKeyboard := gotele.NewReplyKeyboardOneTime(
		gotele.NewKeyboardGrid(
			gotele.NewKeyboardRow(
				gotele.NewKeyboardButton("New Option 1"),
				gotele.NewKeyboardButton("New Option 2"),
			),
		),
	)

	markupOptions := &gotele.EditMessageReplyMarkupOptions{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: newKeyboard,
	}

	err = bot.EditMessageReplyMarkup(markupOptions)
	if err != nil {
		log.Printf("Failed to edit reply markup: %v", err)
	} else {
		fmt.Println("Reply markup edited successfully")
	}

	// Example 5: Advanced keyboard with different button types
	fmt.Println("\n=== Example 5: Advanced Keyboard ===")

	advancedKeyboard := gotele.NewReplyKeyboard(
		gotele.NewKeyboardGrid(
			gotele.NewKeyboardRow(
				gotele.NewKeyboardButton("Regular Button"),
				gotele.NewKeyboardButtonWithRequest("Share Contact", "contact"),
			),
			gotele.NewKeyboardRow(
				gotele.NewKeyboardButtonWithRequest("Share Location", "location"),
				gotele.NewKeyboardButton("Another Button"),
			),
		),
	)

	advancedMessageOptions := &gotele.SendMessageOptions{
		ReplyMarkup: advancedKeyboard,
	}

	err = bot.SendMessageAdvanced(chatID, "This keyboard has different button types:", advancedMessageOptions)
	if err != nil {
		log.Printf("Failed to send advanced message: %v", err)
	} else {
		fmt.Println("Advanced keyboard message sent")
	}

	fmt.Println("\n=== All examples completed ===")
}
