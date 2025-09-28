package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	webhookURL := "https://your-domain.com/webhook" // Replace with your actual domain
	secretToken := "your-secret-token-123"          // Replace with your secret token

	fmt.Println("=== Webhook Examples ===")

	// Example 1: Set webhook
	fmt.Println("\n1. Setting webhook...")
	webhookOptions := &gotele.SetWebhookOptions{
		URL:                webhookURL,
		SecretToken:        secretToken,
		MaxConnections:     40,
		AllowedUpdates:     []string{"message", "callback_query", "inline_query"},
		DropPendingUpdates: true,
	}

	err := bot.SetWebhook(webhookOptions)
	if err != nil {
		fmt.Printf("Error setting webhook: %v\n", err)
	} else {
		fmt.Println("Webhook set successfully!")
	}

	// Example 2: Get webhook info
	fmt.Println("\n2. Getting webhook info...")
	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		fmt.Printf("Error getting webhook info: %v\n", err)
	} else {
		fmt.Printf("Webhook URL: %s\n", webhookInfo.URL)
		fmt.Printf("Pending updates: %d\n", webhookInfo.PendingUpdateCount)
		fmt.Printf("Max connections: %d\n", webhookInfo.MaxConnections)
		fmt.Printf("Allowed updates: %v\n", webhookInfo.AllowedUpdates)
		if webhookInfo.LastErrorMessage != "" {
			fmt.Printf("Last error: %s\n", webhookInfo.LastErrorMessage)
		}
	}

	// Example 3: Simple webhook handler
	fmt.Println("\n3. Creating simple webhook handler...")
	_ = func(update *gotele.Update) error {
		if update.Message != nil {
			fmt.Printf("Received message: %s from user %d\n",
				update.Message.Text,
				update.Message.From.ID)

			// Echo the message back
			replyOptions := &gotele.SendMessageOptions{
				ReplyToMessageID: update.Message.MessageID,
			}

			return bot.SendMessageAdvanced(update.Message.Chat.ID, "Echo: "+update.Message.Text, replyOptions)
		}
		return nil
	}

	// Example 4: Advanced webhook handler with routing
	fmt.Println("\n4. Creating advanced webhook handler...")
	advancedHandler := func(update *gotele.Update) error {
		// Route based on update type
		handlers := map[string]gotele.WebhookHandler{
			"message": func(update *gotele.Update) error {
				msg := update.Message
				fmt.Printf("Message from %s: %s\n", msg.From.FirstName, msg.Text)

				// Handle different message types
				if msg.Text == "/start" {
					keyboard := gotele.InlineKeyboardMarkup{
						InlineKeyboard: [][]gotele.InlineKeyboardButton{
							{
								{Text: "Button 1", CallbackData: "btn1"},
								{Text: "Button 2", CallbackData: "btn2"},
							},
						},
					}

					options := &gotele.SendMessageOptions{
						ReplyMarkup: keyboard,
					}
					return bot.SendMessageAdvanced(msg.Chat.ID, "Welcome! Choose an option:", options)
				}

				return nil
			},
			"callback_query": func(update *gotele.Update) error {
				callback := update.CallbackQuery
				fmt.Printf("Callback query from %s: %s\n",
					callback.From.FirstName,
					callback.Data)

				// Answer callback query
				answerOptions := &gotele.AnswerCallbackQueryOptions{
					CallbackQueryID: callback.ID,
					Text:            "Button clicked!",
					ShowAlert:       false,
				}
				return bot.AnswerCallbackQuery(answerOptions)
			},
			"inline_query": func(update *gotele.Update) error {
				query := update.InlineQuery
				fmt.Printf("Inline query from %s: %s\n",
					query.From.FirstName,
					query.Query)

				// Answer inline query
				results := []gotele.InlineQueryResult{
					{
						Type:  "article",
						ID:    "1",
						Title: "Result 1",
						InputMessageContent: gotele.InputTextMessageContent{
							MessageText: "This is result 1",
						},
					},
				}

				answerOptions := &gotele.AnswerInlineQueryOptions{
					InlineQueryID: query.ID,
					Results:       results,
					CacheTime:     300,
				}
				return bot.AnswerInlineQuery(answerOptions)
			},
			"default": func(update *gotele.Update) error {
				fmt.Printf("Unhandled update type: %d\n", update.UpdateID)
				return nil
			},
		}

		return bot.ProcessWebhookUpdate(update, handlers)
	}

	// Example 5: Start webhook server
	fmt.Println("\n5. Starting webhook server...")
	serverConfig := &gotele.WebhookServer{
		Port:           "8080",
		Path:           "/webhook",
		Handler:        advancedHandler,
		SecretToken:    secretToken,
		AllowedUpdates: []string{"message", "callback_query", "inline_query"},
	}

	// Start server in a goroutine
	go func() {
		fmt.Println("Webhook server starting on port 8080...")
		if err := bot.StartWebhookServer(serverConfig); err != nil {
			fmt.Printf("Webhook server error: %v\n", err)
		}
	}()

	// Example 6: Webhook with middleware
	fmt.Println("\n6. Creating webhook with middleware...")
	mux := http.NewServeMux()

	// Add logging middleware
	loggedHandler := gotele.WebhookLogger(bot.WebhookHandlerFunc(secretToken, advancedHandler))
	mux.HandleFunc("/webhook", loggedHandler.ServeHTTP)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server with middleware
	go func() {
		fmt.Println("Webhook server with middleware starting on port 8081...")
		server := &http.Server{
			Addr:    ":8081",
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Webhook server error: %v\n", err)
		}
	}()

	// Example 7: Webhook with context timeout
	fmt.Println("\n7. Testing webhook with context timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set webhook with context
	err = bot.SetWebhookWithContext(ctx, webhookOptions)
	if err != nil {
		fmt.Printf("Error setting webhook with context: %v\n", err)
	}

	// Example 8: Webhook validation
	fmt.Println("\n8. Testing webhook signature validation...")
	testBody := `{"update_id":123,"message":{"message_id":1,"text":"Hello"}}`

	// Test valid signature
	valid := gotele.ValidateWebhookSignature(secretToken, testBody, "sha256=valid_signature")
	fmt.Printf("Valid signature test: %t\n", valid)

	// Test invalid signature
	invalid := gotele.ValidateWebhookSignature(secretToken, testBody, "invalid")
	fmt.Printf("Invalid signature test: %t\n", invalid)

	// Example 9: Delete webhook
	fmt.Println("\n9. Deleting webhook...")
	err = bot.DeleteWebhook(true) // Drop pending updates
	if err != nil {
		fmt.Printf("Error deleting webhook: %v\n", err)
	} else {
		fmt.Println("Webhook deleted successfully!")
	}

	// Example 10: Webhook with custom certificate
	fmt.Println("\n10. Setting webhook with custom certificate...")
	certOptions := &gotele.SetWebhookOptions{
		URL: webhookURL,
		Certificate: gotele.InputFile{
			FilePath: "certificate.pem", // Replace with actual certificate file
		},
		SecretToken: secretToken,
	}

	err = bot.SetWebhook(certOptions)
	if err != nil {
		fmt.Printf("Error setting webhook with certificate: %v\n", err)
	} else {
		fmt.Println("Webhook with certificate set successfully!")
	}

	// Keep the program running to demonstrate webhook server
	fmt.Println("\n=== Webhook servers are running ===")
	fmt.Println("Send messages to your bot to test the webhook handlers")
	fmt.Println("Press Ctrl+C to stop")

	// Wait for interrupt
	select {}
}
