package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gotele "github.com/repute-software/gotele/telegram"
)

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.Index(line, "="); i >= 0 {
			key := strings.TrimSpace(line[:i])
			val := strings.TrimSpace(line[i+1:])
			val = strings.Trim(val, `"'`)
			if key != "" {
				if _, exists := os.LookupEnv(key); !exists {
					_ = os.Setenv(key, val)
				}
			}
		}
	}
}

func main() {
	// Load .env for local development (no-op if missing)
	loadEnvFile(".env")

	token := os.Getenv("BOT_DEV_TOKEN")
	if token == "" {
		log.Fatal("BOT_DEV_TOKEN is not set")
	}

	// Create bot with custom timeout
	bot := gotele.NewBotWithTimeout(token, 15*time.Second)
	fmt.Println("Bot created with 15 second timeout")

	// Example 1: Using context with timeout
	fmt.Println("\n=== Example 1: Context with Timeout ===")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updates, err := bot.GetUpdatesWithContext(ctx, 0)
	if err != nil {
		log.Printf("GetUpdates failed: %v", err)
	} else {
		fmt.Printf("Received %d pending update(s) with context.\n", len(updates))
	}

	// Example 2: Using the old method (backward compatibility)
	fmt.Println("\n=== Example 2: Backward Compatibility ===")
	updates2, err := bot.GetUpdates(0)
	if err != nil {
		log.Printf("GetUpdates (old method) failed: %v", err)
	} else {
		fmt.Printf("Received %d pending update(s) with old method.\n", len(updates2))
	}

	// Example 3: Context cancellation
	fmt.Println("\n=== Example 3: Context Cancellation ===")
	ctx2, cancel2 := context.WithCancel(context.Background())

	// Cancel after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		cancel2()
		fmt.Println("Context cancelled after 1 second")
	}()

	updates3, err := bot.GetUpdatesWithContext(ctx2, 0)
	if err != nil {
		fmt.Printf("GetUpdates with cancellation: %v\n", err)
	} else {
		fmt.Printf("Received %d updates before cancellation\n", len(updates3))
	}

	fmt.Println("\n=== Bot ready with context support! ===")
}
