package gotele

import (
	"testing"
)

func TestAPIError(t *testing.T) {
	err := &APIError{
		ErrorCode:   400,
		Description: "Bad Request: chat not found",
	}

	expected := "telegram API error 400: Bad Request: chat not found"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	if err.IsRetryable() {
		t.Error("Expected error 400 to not be retryable")
	}
}

func TestAPIErrorRetryable(t *testing.T) {
	err := &APIError{
		ErrorCode:   429,
		Description: "Too Many Requests: retry after 60",
	}

	if !err.IsRetryable() {
		t.Error("Expected error 429 to be retryable")
	}
}

func TestHTTPError(t *testing.T) {
	err := &HTTPError{
		StatusCode: 500,
		Status:     "Internal Server Error",
		Body:       "Server is down",
	}

	expected := "HTTP 500 Internal Server Error: Server is down"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	if !err.IsRetryable() {
		t.Error("Expected HTTP 500 to be retryable")
	}
}

func TestAPIResponseToError(t *testing.T) {
	// Test successful response
	resp := &APIResponse{
		Ok:     true,
		Result: "success",
	}

	if err := resp.ToError(); err != nil {
		t.Errorf("Expected no error for successful response, got %v", err)
	}

	// Test error response
	resp = &APIResponse{
		Ok:          false,
		ErrorCode:   400,
		Description: "Bad Request",
		Parameters:  map[string]interface{}{"retry_after": 60},
	}

	err := resp.ToError()
	if err == nil {
		t.Error("Expected error for failed response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("Expected APIError, got %T", err)
	}

	if apiErr.ErrorCode != 400 {
		t.Errorf("Expected error code 400, got %d", apiErr.ErrorCode)
	}

	if apiErr.Description != "Bad Request" {
		t.Errorf("Expected description 'Bad Request', got %q", apiErr.Description)
	}
}
