package gotele

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebhookInfo(t *testing.T) {
	webhookInfo := WebhookInfo{
		URL:                  "https://example.com/webhook",
		HasCustomCertificate: false,
		PendingUpdateCount:   5,
		IPAddress:            "192.168.1.1",
		LastErrorDate:        1640995200,
		LastErrorMessage:     "Connection timeout",
		MaxConnections:       40,
		AllowedUpdates:       []string{"message", "callback_query"},
	}

	if webhookInfo.URL != "https://example.com/webhook" {
		t.Errorf("Expected URL 'https://example.com/webhook', got %s", webhookInfo.URL)
	}
	if webhookInfo.PendingUpdateCount != 5 {
		t.Errorf("Expected PendingUpdateCount 5, got %d", webhookInfo.PendingUpdateCount)
	}
	if webhookInfo.IPAddress != "192.168.1.1" {
		t.Errorf("Expected IPAddress '192.168.1.1', got %s", webhookInfo.IPAddress)
	}
	if len(webhookInfo.AllowedUpdates) != 2 {
		t.Errorf("Expected 2 allowed updates, got %d", len(webhookInfo.AllowedUpdates))
	}
}

func TestSetWebhookOptions(t *testing.T) {
	options := SetWebhookOptions{
		URL:                "https://example.com/webhook",
		IPAddress:          "192.168.1.1",
		MaxConnections:     40,
		AllowedUpdates:     []string{"message", "callback_query"},
		DropPendingUpdates: true,
		SecretToken:        "secret123",
	}

	if options.URL != "https://example.com/webhook" {
		t.Errorf("Expected URL 'https://example.com/webhook', got %s", options.URL)
	}
	if options.MaxConnections != 40 {
		t.Errorf("Expected MaxConnections 40, got %d", options.MaxConnections)
	}
	if !options.DropPendingUpdates {
		t.Error("Expected DropPendingUpdates to be true")
	}
	if options.SecretToken != "secret123" {
		t.Errorf("Expected SecretToken 'secret123', got %s", options.SecretToken)
	}
}

func TestWebhookServer(t *testing.T) {
	server := WebhookServer{
		Port:           "8080",
		Path:           "/webhook",
		Handler:        func(update *Update) error { return nil },
		SecretToken:    "secret123",
		AllowedUpdates: []string{"message"},
	}

	if server.Port != "8080" {
		t.Errorf("Expected Port '8080', got %s", server.Port)
	}
	if server.Path != "/webhook" {
		t.Errorf("Expected Path '/webhook', got %s", server.Path)
	}
	if server.SecretToken != "secret123" {
		t.Errorf("Expected SecretToken 'secret123', got %s", server.SecretToken)
	}
}

func TestWebhookUpdate(t *testing.T) {
	update := &Update{
		UpdateID: 123,
		Message: &Message{
			MessageID: 1,
			Text:      "Hello",
		},
	}

	webhookUpdate := WebhookUpdate{
		Update:      update,
		Timestamp:   1640995200,
		IPAddress:   "192.168.1.1",
		UserAgent:   "TelegramBot/1.0",
		SecretToken: "secret123",
	}

	if webhookUpdate.Update.UpdateID != 123 {
		t.Errorf("Expected UpdateID 123, got %d", webhookUpdate.Update.UpdateID)
	}
	if webhookUpdate.IPAddress != "192.168.1.1" {
		t.Errorf("Expected IPAddress '192.168.1.1', got %s", webhookUpdate.IPAddress)
	}
	if webhookUpdate.UserAgent != "TelegramBot/1.0" {
		t.Errorf("Expected UserAgent 'TelegramBot/1.0', got %s", webhookUpdate.UserAgent)
	}
}

func TestValidateWebhookSignature(t *testing.T) {
	secretToken := "secret123"
	body := `{"update_id":123,"message":{"message_id":1,"text":"Hello"}}`

	// Create valid signature
	mac := hmac.New(sha256.New, []byte(secretToken))
	mac.Write([]byte(body))
	validSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// Test valid signature
	if !ValidateWebhookSignature(secretToken, body, validSignature) {
		t.Error("Expected valid signature to pass validation")
	}

	// Test invalid signature
	invalidSignature := "sha256=invalid"
	if ValidateWebhookSignature(secretToken, body, invalidSignature) {
		t.Error("Expected invalid signature to fail validation")
	}

	// Test without sha256= prefix
	signatureWithoutPrefix := hex.EncodeToString(mac.Sum(nil))
	if !ValidateWebhookSignature(secretToken, body, signatureWithoutPrefix) {
		t.Error("Expected signature without prefix to pass validation")
	}

	// Test empty secret token
	if ValidateWebhookSignature("", body, validSignature) {
		t.Error("Expected empty secret token to fail validation")
	}

	// Test empty signature
	if ValidateWebhookSignature(secretToken, body, "") {
		t.Error("Expected empty signature to fail validation")
	}
}

func TestWebhookHandlerFunc(t *testing.T) {
	bot := NewBot("test_token")
	secretToken := "secret123"

	// Create test handler
	handlerCalled := false
	handler := func(update *Update) error {
		handlerCalled = true
		if update.UpdateID != 123 {
			t.Errorf("Expected UpdateID 123, got %d", update.UpdateID)
		}
		return nil
	}

	// Create HTTP handler
	httpHandler := bot.WebhookHandlerFunc(secretToken, handler)

	// Create test update
	update := Update{
		UpdateID: 123,
		Message: &Message{
			MessageID: 1,
			Text:      "Hello",
		},
	}

	updateJSON, _ := json.Marshal(update)

	// Create valid signature
	mac := hmac.New(sha256.New, []byte(secretToken))
	mac.Write(updateJSON)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// Test valid request
	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(string(updateJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", signature)

	w := httptest.NewRecorder()
	httpHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	// Test invalid signature
	handlerCalled = false
	req = httptest.NewRequest("POST", "/webhook", strings.NewReader(string(updateJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "invalid")

	w = httptest.NewRecorder()
	httpHandler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	if handlerCalled {
		t.Error("Expected handler not to be called with invalid signature")
	}

	// Test OPTIONS request
	req = httptest.NewRequest("OPTIONS", "/webhook", nil)
	w = httptest.NewRecorder()
	httpHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	// Test non-POST request
	req = httptest.NewRequest("GET", "/webhook", nil)
	w = httptest.NewRecorder()
	httpHandler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestProcessWebhookUpdate(t *testing.T) {
	bot := NewBot("test_token")

	// Create handlers
	messageHandlerCalled := false
	callbackHandlerCalled := false
	defaultHandlerCalled := false

	handlers := map[string]WebhookHandler{
		"message": func(update *Update) error {
			messageHandlerCalled = true
			return nil
		},
		"callback_query": func(update *Update) error {
			callbackHandlerCalled = true
			return nil
		},
		"default": func(update *Update) error {
			defaultHandlerCalled = true
			return nil
		},
	}

	// Test message update
	messageUpdate := &Update{
		UpdateID: 123,
		Message: &Message{
			MessageID: 1,
			Text:      "Hello",
		},
	}

	err := bot.ProcessWebhookUpdate(messageUpdate, handlers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !messageHandlerCalled {
		t.Error("Expected message handler to be called")
	}

	// Test callback query update
	messageHandlerCalled = false
	callbackUpdate := &Update{
		UpdateID:      124,
		CallbackQuery: &CallbackQuery{},
	}

	err = bot.ProcessWebhookUpdate(callbackUpdate, handlers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !callbackHandlerCalled {
		t.Error("Expected callback handler to be called")
	}

	// Test unknown update type
	messageHandlerCalled = false
	callbackHandlerCalled = false
	unknownUpdate := &Update{
		UpdateID: 125,
	}

	err = bot.ProcessWebhookUpdate(unknownUpdate, handlers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !defaultHandlerCalled {
		t.Error("Expected default handler to be called")
	}

	// Test with no handlers
	handlers = map[string]WebhookHandler{}
	err = bot.ProcessWebhookUpdate(messageUpdate, handlers)
	if err == nil {
		t.Error("Expected error for unknown update type with no handlers")
	}
}

func TestGetClientIP(t *testing.T) {
	// Test X-Forwarded-For header
	req := httptest.NewRequest("POST", "/webhook", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")

	ip := getClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("Expected IP '192.168.1.1', got %s", ip)
	}

	// Test X-Real-IP header
	req = httptest.NewRequest("POST", "/webhook", nil)
	req.Header.Set("X-Real-IP", "192.168.1.2")

	ip = getClientIP(req)
	if ip != "192.168.1.2" {
		t.Errorf("Expected IP '192.168.1.2', got %s", ip)
	}

	// Test RemoteAddr fallback
	req = httptest.NewRequest("POST", "/webhook", nil)
	req.RemoteAddr = "192.168.1.3:12345"

	ip = getClientIP(req)
	if ip != "192.168.1.3" {
		t.Errorf("Expected IP '192.168.1.3', got %s", ip)
	}
}

func TestWebhookWithContext(t *testing.T) {
	bot := NewBot("test_token")

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	options := &SetWebhookOptions{
		URL: "https://example.com/webhook",
	}

	err := bot.SetWebhookWithContext(ctx, options)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	// Test with timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel2()

	// Wait for timeout
	time.Sleep(1 * time.Millisecond)

	err = bot.SetWebhookWithContext(ctx2, options)
	if err == nil {
		t.Error("Expected error for timed out context")
	}
}

func TestWebhookMiddleware(t *testing.T) {
	secretToken := "secret123"

	// Create next handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create middleware
	middleware := WebhookMiddleware(secretToken, nextHandler)

	// Test valid request
	body := `{"test": "data"}`
	mac := hmac.New(sha256.New, []byte(secretToken))
	mac.Write([]byte(body))
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", signature)

	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test invalid signature
	req = httptest.NewRequest("POST", "/webhook", strings.NewReader(body))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "invalid")

	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	// Test OPTIONS request
	req = httptest.NewRequest("OPTIONS", "/webhook", nil)
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}
}

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	// Test WriteHeader
	rw.WriteHeader(http.StatusCreated)
	if rw.statusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rw.statusCode)
	}

	// Test Write
	data := []byte("test data")
	n, err := rw.Write(data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected %d bytes written, got %d", len(data), n)
	}
}
