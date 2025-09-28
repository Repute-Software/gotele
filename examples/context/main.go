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

	// Example 1: Basic context usage with timeout
	fmt.Println("=== Example 1: Context with Timeout ===")
	bot := gotele.NewBot(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := bot.SendMessageWithContext(ctx, 123456789, "Hello with context!")
	if err != nil {
		fmt.Printf("SendMessage error: %v\n", err)
	} else {
		fmt.Println("Message sent successfully!")
	}

	// Example 2: Context cancellation
	fmt.Println("\n=== Example 2: Context Cancellation ===")
	ctx2, cancel2 := context.WithCancel(context.Background())

	// Cancel after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		cancel2()
		fmt.Println("Context cancelled!")
	}()

	err = bot.SendMessageWithContext(ctx2, 123456789, "This might be cancelled")
	if err != nil {
		fmt.Printf("SendMessage error (expected due to cancellation): %v\n", err)
	}

	// Example 3: Custom timeout bot
	fmt.Println("\n=== Example 3: Custom Timeout Bot ===")
	fastBot := gotele.NewBotWithTimeout(token, 2*time.Second)

	err = fastBot.SendMessage(123456789, "Quick message with 2s timeout")
	if err != nil {
		fmt.Printf("Fast bot error: %v\n", err)
	}

	// Example 4: Long polling with context
	fmt.Println("\n=== Example 4: GetUpdates with Context ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()

	updates, err := bot.GetUpdatesWithContext(ctx3, 0)
	if err != nil {
		fmt.Printf("GetUpdates error: %v\n", err)
	} else {
		fmt.Printf("Received %d updates\n", len(updates))
	}

	// Example 5: Context with deadline
	fmt.Println("\n=== Example 5: Context with Deadline ===")
	deadline := time.Now().Add(3 * time.Second)
	ctx4, cancel4 := context.WithDeadline(context.Background(), deadline)
	defer cancel4()

	err = bot.SendMessageWithContext(ctx4, 123456789, "Message with deadline")
	if err != nil {
		fmt.Printf("SendMessage with deadline error: %v\n", err)
	}

	// Example 6: Context with values
	fmt.Println("\n=== Example 6: Context with Values ===")
	ctx5 := context.WithValue(context.Background(), "request_id", "req_123")
	ctx5 = context.WithValue(ctx5, "user_id", "user_456")

	err = bot.SendMessageWithContext(ctx5, 123456789, "Message with context values")
	if err != nil {
		fmt.Printf("SendMessage with values error: %v\n", err)
	}

	fmt.Println("\n=== All examples completed ===")
}
