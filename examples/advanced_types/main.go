package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gotele "github.com/repute-software/gotele/telegram"
)

func main() {
	// Load environment variables
	token := os.Getenv("BOT_DEV_TOKEN")
	if token == "" {
		log.Fatal("BOT_DEV_TOKEN is not set")
	}

	bot := gotele.NewBot(token)
	chatID := int64(123456789) // Replace with actual chat ID

	fmt.Println("=== Advanced Types Examples ===")

	// Example 1: Send message with inline keyboard
	fmt.Println("\n1. Sending message with inline keyboard...")
	keyboard := gotele.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotele.InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "btn1"},
				{Text: "Button 2", CallbackData: "btn2"},
			},
			{
				{Text: "Visit Website", URL: "https://example.com"},
			},
		},
	}

	options := &gotele.SendMessageOptions{
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}

	err := bot.SendMessageAdvanced(chatID, "*Hello!* Choose an option:", options)
	if err != nil {
		fmt.Printf("Error sending message with keyboard: %v\n", err)
	}

	// Example 2: Send message with reply keyboard
	fmt.Println("\n2. Sending message with reply keyboard...")
	replyKeyboard := gotele.ReplyKeyboardMarkup{
		Keyboard: [][]gotele.KeyboardButton{
			{
				{Text: "Option 1"},
				{Text: "Option 2"},
			},
			{
				{Text: "Share Contact", RequestContact: true},
				{Text: "Share Location", RequestLocation: true},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	options2 := &gotele.SendMessageOptions{
		ReplyMarkup: replyKeyboard,
	}

	err = bot.SendMessageAdvanced(chatID, "Please choose an option:", options2)
	if err != nil {
		fmt.Printf("Error sending message with reply keyboard: %v\n", err)
	}

	// Example 3: Send photo with caption
	fmt.Println("\n3. Sending photo with caption...")
	photoOptions := &gotele.SendPhotoOptions{
		ChatID:    chatID,
		Photo:     "https://picsum.photos/400/300", // Random image URL
		Caption:   "Here's a random image! 📸",
		ParseMode: "Markdown",
	}

	err = bot.SendPhoto(photoOptions)
	if err != nil {
		fmt.Printf("Error sending photo: %v\n", err)
	}

	// Example 4: Send message with entities
	fmt.Println("\n4. Sending message with entities...")
	entities := []gotele.MessageEntity{
		{
			Type:   "bold",
			Offset: 0,
			Length: 5,
		},
		{
			Type:   "italic",
			Offset: 6,
			Length: 4,
		},
		{
			Type:   "code",
			Offset: 11,
			Length: 8,
		},
	}

	options3 := &gotele.SendMessageOptions{
		Entities: entities,
	}

	err = bot.SendMessageAdvanced(chatID, "Bold text and italic and `code here`", options3)
	if err != nil {
		fmt.Printf("Error sending message with entities: %v\n", err)
	}

	// Example 5: Send message with context and timeout
	fmt.Println("\n5. Sending message with context timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options4 := &gotele.SendMessageOptions{
		ParseMode:           "HTML",
		DisableNotification: true,
	}

	err = bot.SendMessageAdvancedWithContext(ctx, chatID, "<b>Bold text</b> with <i>HTML formatting</i>", options4)
	if err != nil {
		fmt.Printf("Error sending message with context: %v\n", err)
	}

	// Example 6: Edit message
	fmt.Println("\n6. Editing a message...")
	// First send a message
	err = bot.SendMessage(chatID, "This message will be edited")
	if err != nil {
		fmt.Printf("Error sending initial message: %v\n", err)
	} else {
		// Wait a bit then edit (in real usage, you'd get the message ID from the response)
		time.Sleep(1 * time.Second)

		editOptions := &gotele.EditMessageTextOptions{
			ChatID:    chatID,
			MessageID: 1, // This would be the actual message ID from the response
			Text:      "This message has been edited! ✏️",
			ParseMode: "Markdown",
		}

		err = bot.EditMessageText(editOptions)
		if err != nil {
			fmt.Printf("Error editing message: %v\n", err)
		}
	}

	// Example 7: Send message with protection and reply
	fmt.Println("\n7. Sending protected message with reply...")
	options5 := &gotele.SendMessageOptions{
		ProtectContent:   true,
		ReplyToMessageID: 1, // This would be the actual message ID
		ParseMode:        "Markdown",
	}

	err = bot.SendMessageAdvanced(chatID, "This is a *protected* message that replies to another message", options5)
	if err != nil {
		fmt.Printf("Error sending protected message: %v\n", err)
	}

	fmt.Println("\n=== All examples completed ===")
}
