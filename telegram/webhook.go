package gotele

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SetWebhook sets a webhook for the bot
func (b *Bot) SetWebhook(options *SetWebhookOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SetWebhookWithContext(ctx, options)
}

// SetWebhookWithContext sets a webhook for the bot with context support
func (b *Bot) SetWebhookWithContext(ctx context.Context, options *SetWebhookOptions) error {
	reqBody := map[string]interface{}{
		"url": options.URL,
	}

	if options.Certificate.FilePath != "" || options.Certificate.URL != "" || len(options.Certificate.Data) > 0 {
		// Handle certificate upload
		certUpload, err := b.prepareFileUpload(options.Certificate, "certificate")
		if err != nil {
			return fmt.Errorf("failed to prepare certificate: %w", err)
		}

		fields := map[string]string{
			"url": options.URL,
		}

		if options.IPAddress != "" {
			fields["ip_address"] = options.IPAddress
		}
		if options.MaxConnections > 0 {
			fields["max_connections"] = fmt.Sprintf("%d", options.MaxConnections)
		}
		if len(options.AllowedUpdates) > 0 {
			allowedUpdatesJSON, err := json.Marshal(options.AllowedUpdates)
			if err != nil {
				return fmt.Errorf("failed to marshal allowed updates: %w", err)
			}
			fields["allowed_updates"] = string(allowedUpdatesJSON)
		}
		if options.DropPendingUpdates {
			fields["drop_pending_updates"] = "true"
		}
		if options.SecretToken != "" {
			fields["secret_token"] = options.SecretToken
		}

		_, err = b.makeMultipartRequest(ctx, "/setWebhook", fields, []FileUpload{certUpload})
		return err
	}

	// Regular JSON request without certificate
	if options.IPAddress != "" {
		reqBody["ip_address"] = options.IPAddress
	}
	if options.MaxConnections > 0 {
		reqBody["max_connections"] = options.MaxConnections
	}
	if len(options.AllowedUpdates) > 0 {
		reqBody["allowed_updates"] = options.AllowedUpdates
	}
	if options.DropPendingUpdates {
		reqBody["drop_pending_updates"] = true
	}
	if options.SecretToken != "" {
		reqBody["secret_token"] = options.SecretToken
	}

	_, err := b.makeRequest(ctx, "POST", "/setWebhook", reqBody)
	return err
}

// DeleteWebhook removes the webhook for the bot
func (b *Bot) DeleteWebhook(dropPendingUpdates bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.DeleteWebhookWithContext(ctx, dropPendingUpdates)
}

// DeleteWebhookWithContext removes the webhook for the bot with context support
func (b *Bot) DeleteWebhookWithContext(ctx context.Context, dropPendingUpdates bool) error {
	reqBody := map[string]interface{}{}
	if dropPendingUpdates {
		reqBody["drop_pending_updates"] = true
	}

	_, err := b.makeRequest(ctx, "POST", "/deleteWebhook", reqBody)
	return err
}

// GetWebhookInfo gets information about the current webhook
func (b *Bot) GetWebhookInfo() (*WebhookInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.GetWebhookInfoWithContext(ctx)
}

// GetWebhookInfoWithContext gets information about the current webhook with context support
func (b *Bot) GetWebhookInfoWithContext(ctx context.Context) (*WebhookInfo, error) {
	resp, err := b.makeRequest(ctx, "GET", "/getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	// Parse the result into WebhookInfo
	var webhookInfo WebhookInfo
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := json.Unmarshal(resultBytes, &webhookInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook info: %w", err)
	}

	return &webhookInfo, nil
}

// ValidateWebhookSignature validates the webhook signature for security
func ValidateWebhookSignature(secretToken, body, signature string) bool {
	if secretToken == "" || signature == "" {
		return false
	}

	// Remove "sha256=" prefix if present
	if strings.HasPrefix(signature, "sha256=") {
		signature = signature[7:]
	}

	// Create HMAC
	mac := hmac.New(sha256.New, []byte(secretToken))
	mac.Write([]byte(body))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// WebhookHandlerFunc creates an HTTP handler for webhook updates
func (b *Bot) WebhookHandlerFunc(secretToken string, handler WebhookHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Telegram-Bot-Api-Secret-Token")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only accept POST requests
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Validate signature if secret token is provided
		if secretToken != "" {
			signature := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if !ValidateWebhookSignature(secretToken, string(body), signature) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}
		}

		// Parse update
		var update Update
		if err := json.Unmarshal(body, &update); err != nil {
			http.Error(w, "Failed to parse update", http.StatusBadRequest)
			return
		}

		// Create webhook update with metadata (for future use)
		_ = &WebhookUpdate{
			Update:      &update,
			Timestamp:   time.Now().Unix(),
			IPAddress:   getClientIP(r),
			UserAgent:   r.Header.Get("User-Agent"),
			SecretToken: secretToken,
		}

		// Call handler
		if err := handler(&update); err != nil {
			http.Error(w, "Handler error", http.StatusInternalServerError)
			return
		}

		// Send success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// StartWebhookServer starts an HTTP server for webhook updates
func (b *Bot) StartWebhookServer(config *WebhookServer) error {
	// Create HTTP handler
	handler := b.WebhookHandlerFunc(config.SecretToken, config.Handler)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: http.HandlerFunc(handler),
	}

	// Start server
	return server.ListenAndServe()
}

// StartWebhookServerTLS starts an HTTPS server for webhook updates
func (b *Bot) StartWebhookServerTLS(config *WebhookServer, certFile, keyFile string) error {
	// Create HTTP handler
	handler := b.WebhookHandlerFunc(config.SecretToken, config.Handler)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: http.HandlerFunc(handler),
	}

	// Start server with TLS
	return server.ListenAndServeTLS(certFile, keyFile)
}

// ProcessWebhookUpdate processes a webhook update with routing
func (b *Bot) ProcessWebhookUpdate(update *Update, handlers map[string]WebhookHandler) error {
	// Route based on update type
	var handlerType string

	if update.Message != nil {
		handlerType = "message"
	} else if update.EditedMessage != nil {
		handlerType = "edited_message"
	} else if update.ChannelPost != nil {
		handlerType = "channel_post"
	} else if update.EditedChannelPost != nil {
		handlerType = "edited_channel_post"
	} else if update.InlineQuery != nil {
		handlerType = "inline_query"
	} else if update.ChosenInlineResult != nil {
		handlerType = "chosen_inline_result"
	} else if update.CallbackQuery != nil {
		handlerType = "callback_query"
	} else if update.ShippingQuery != nil {
		handlerType = "shipping_query"
	} else if update.PreCheckoutQuery != nil {
		handlerType = "pre_checkout_query"
	} else if update.Poll != nil {
		handlerType = "poll"
	} else if update.PollAnswer != nil {
		handlerType = "poll_answer"
	} else if update.MyChatMember != nil {
		handlerType = "my_chat_member"
	} else if update.ChatMember != nil {
		handlerType = "chat_member"
	} else if update.ChatJoinRequest != nil {
		handlerType = "chat_join_request"
	}

	// Call specific handler if available
	if handler, exists := handlers[handlerType]; exists {
		return handler(update)
	}

	// Call default handler if available
	if defaultHandler, exists := handlers["default"]; exists {
		return defaultHandler(update)
	}

	// No handler found
	return fmt.Errorf("no handler found for update type: %s", handlerType)
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// WebhookMiddleware creates middleware for webhook processing
func WebhookMiddleware(secretToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Telegram-Bot-Api-Secret-Token")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only accept POST requests
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Validate signature if secret token is provided
		if secretToken != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			// Restore body for next handler
			r.Body = io.NopCloser(strings.NewReader(string(body)))

			signature := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if !ValidateWebhookSignature(secretToken, string(body), signature) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}
		}

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// WebhookLogger creates a logging middleware for webhook requests
func WebhookLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log request
		duration := time.Since(start)
		fmt.Printf("[WEBHOOK] %s %s %d %v %s\n",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			getClientIP(r))
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
