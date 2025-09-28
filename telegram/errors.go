package gotele

import (
	"fmt"
)

// APIError represents a Telegram Bot API error response
type APIError struct {
	ErrorCode   int                    `json:"error_code"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("telegram API error %d: %s", e.ErrorCode, e.Description)
}

// IsRetryable returns true if the error is retryable
func (e *APIError) IsRetryable() bool {
	// Common retryable error codes
	retryableCodes := map[int]bool{
		429: true, // Too Many Requests
		500: true, // Internal Server Error
		502: true, // Bad Gateway
		503: true, // Service Unavailable
		504: true, // Gateway Timeout
	}
	return retryableCodes[e.ErrorCode]
}

// APIResponse represents the standard Telegram Bot API response format
type APIResponse struct {
	Ok          bool        `json:"ok"`
	Result      interface{} `json:"result,omitempty"`
	ErrorCode   int         `json:"error_code,omitempty"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
}

// ToError converts an APIResponse to an APIError if the response indicates an error
func (r *APIResponse) ToError() error {
	if r.Ok {
		return nil
	}

	return &APIError{
		ErrorCode:   r.ErrorCode,
		Description: r.Description,
		Parameters:  r.Parameters.(map[string]interface{}),
	}
}

// HTTPError represents HTTP-level errors
type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d %s: %s", e.StatusCode, e.Status, e.Body)
}

// IsRetryable returns true if the HTTP error is retryable
func (e *HTTPError) IsRetryable() bool {
	return e.StatusCode >= 500 || e.StatusCode == 429
}
