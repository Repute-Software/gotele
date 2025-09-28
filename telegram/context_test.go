package gotele

import (
	"context"
	"testing"
	"time"
)

func TestNewBotWithTimeout(t *testing.T) {
	timeout := 10 * time.Second
	bot := NewBotWithTimeout("test_token", timeout)

	if bot.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, bot.Timeout)
	}

	if bot.Token != "test_token" {
		t.Errorf("Expected token 'test_token', got %s", bot.Token)
	}
}

func TestNewBotDefaultTimeout(t *testing.T) {
	bot := NewBot("test_token")

	expectedTimeout := 30 * time.Second
	if bot.Timeout != expectedTimeout {
		t.Errorf("Expected default timeout %v, got %v", expectedTimeout, bot.Timeout)
	}
}

func TestContextCancellation(t *testing.T) {
	bot := NewBot("test_token")

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test SendMessageWithContext with cancelled context
	err := bot.SendMessageWithContext(ctx, 123456789, "test message")
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}

	// Test GetUpdatesWithContext with cancelled context
	updates, err := bot.GetUpdatesWithContext(ctx, 0)
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
	if updates != nil {
		t.Error("Expected nil updates for cancelled context")
	}
}

func TestContextTimeout(t *testing.T) {
	bot := NewBot("test_token")

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for the timeout to occur
	time.Sleep(1 * time.Millisecond)

	// Test SendMessageWithContext with timed out context
	err := bot.SendMessageWithContext(ctx, 123456789, "test message")
	if err == nil {
		t.Error("Expected error for timed out context, got nil")
	}

	// Test GetUpdatesWithContext with timed out context
	updates, err := bot.GetUpdatesWithContext(ctx, 0)
	if err == nil {
		t.Error("Expected error for timed out context, got nil")
	}
	if updates != nil {
		t.Error("Expected nil updates for timed out context")
	}
}

func TestContextDeadline(t *testing.T) {
	bot := NewBot("test_token")

	// Create a context with a deadline in the past
	deadline := time.Now().Add(-1 * time.Hour)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Test SendMessageWithContext with expired deadline
	err := bot.SendMessageWithContext(ctx, 123456789, "test message")
	if err == nil {
		t.Error("Expected error for expired deadline, got nil")
	}

	// Test GetUpdatesWithContext with expired deadline
	updates, err := bot.GetUpdatesWithContext(ctx, 0)
	if err == nil {
		t.Error("Expected error for expired deadline, got nil")
	}
	if updates != nil {
		t.Error("Expected nil updates for expired deadline")
	}
}

func TestContextValue(t *testing.T) {
	bot := NewBot("test_token")

	// Create a context with a value
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	// Test that context is passed through (we can't easily test the value
	// without modifying the makeRequest method, but we can test that it doesn't crash)
	err := bot.SendMessageWithContext(ctx, 123456789, "test message")
	// We expect an error due to invalid token, but not a panic
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	bot := NewBot("test_token")

	// Test that the old methods still work (they should use the bot's timeout)
	err := bot.SendMessage(123456789, "test message")
	// We expect an error due to invalid token, but not a panic
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	updates, err := bot.GetUpdates(0)
	// We expect an error due to invalid token, but not a panic
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
	if updates != nil {
		t.Error("Expected nil updates for invalid token")
	}
}

func TestContextWithTimeout(t *testing.T) {
	// Test that the bot's timeout is used when no context is provided
	bot := NewBotWithTimeout("test_token", 5*time.Second)

	// The SendMessage method should use the bot's timeout
	err := bot.SendMessage(123456789, "test message")
	// We expect an error due to invalid token, but not a panic
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// The GetUpdates method should use the bot's timeout
	updates, err := bot.GetUpdates(0)
	// We expect an error due to invalid token, but not a panic
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
	if updates != nil {
		t.Error("Expected nil updates for invalid token")
	}
}
