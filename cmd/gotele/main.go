package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

	bot := gotele.NewBot(token)

	updates, err := bot.GetUpdates(0)
	if err != nil {
		log.Fatalf("GetUpdates failed: %v", err)
	}

	fmt.Printf("Bot ready. Received %d pending update(s).\n", len(updates))
}
