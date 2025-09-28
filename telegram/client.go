package gotele

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NewBot creates a new Bot instance
func NewBot(token string) *Bot {
	return &Bot{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Client:  &http.Client{},
		Timeout: 30 * time.Second, // Default 30 second timeout
	}
}

// NewBotWithTimeout creates a new Bot instance with custom timeout
func NewBotWithTimeout(token string, timeout time.Duration) *Bot {
	return &Bot{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Client:  &http.Client{},
		Timeout: timeout,
	}
}

// makeRequest makes an HTTP request to the Telegram API and handles the response
func (b *Bot) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequestWithContext(ctx, "GET", b.BaseURL+endpoint, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, b.BaseURL+endpoint, reqBody)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}
	}

	// Parse API response
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API-level errors
	if err := apiResp.ToError(); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

type sendMessageRequest struct {
	ChatID                   int64           `json:"chat_id"`
	MessageThreadID          int             `json:"message_thread_id,omitempty"`
	Text                     string          `json:"text,omitempty"`
	ParseMode                string          `json:"parse_mode,omitempty"`
	Entities                 []MessageEntity `json:"entities,omitempty"`
	DisableWebPagePreview    bool            `json:"disable_web_page_preview,omitempty"`
	DisableNotification      bool            `json:"disable_notification,omitempty"`
	ProtectContent           bool            `json:"protect_content,omitempty"`
	ReplyToMessageID         int             `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply,omitempty"`
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`
}

// SendMessage sends a message to a specific chat
func (b *Bot) SendMessage(chatID int64, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendMessageWithContext(ctx, chatID, text)
}

// SendMessageWithContext sends a message to a specific chat with context support
func (b *Bot) SendMessageWithContext(ctx context.Context, chatID int64, text string) error {
	reqBody := sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	_, err := b.makeRequest(ctx, "POST", "/sendMessage", reqBody)
	return err
}

// SendMessageOptions represents options for sending a message
type SendMessageOptions struct {
	MessageThreadID          int
	ParseMode                string
	Entities                 []MessageEntity
	DisableWebPagePreview    bool
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

// SendMessageAdvanced sends a message with advanced options
func (b *Bot) SendMessageAdvanced(chatID int64, text string, options *SendMessageOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendMessageAdvancedWithContext(ctx, chatID, text, options)
}

// SendMessageAdvancedWithContext sends a message with advanced options and context support
func (b *Bot) SendMessageAdvancedWithContext(ctx context.Context, chatID int64, text string, options *SendMessageOptions) error {
	reqBody := sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	if options != nil {
		reqBody.MessageThreadID = options.MessageThreadID
		reqBody.ParseMode = options.ParseMode
		reqBody.Entities = options.Entities
		reqBody.DisableWebPagePreview = options.DisableWebPagePreview
		reqBody.DisableNotification = options.DisableNotification
		reqBody.ProtectContent = options.ProtectContent
		reqBody.ReplyToMessageID = options.ReplyToMessageID
		reqBody.AllowSendingWithoutReply = options.AllowSendingWithoutReply
		reqBody.ReplyMarkup = options.ReplyMarkup
	}

	_, err := b.makeRequest(ctx, "POST", "/sendMessage", reqBody)
	return err
}

// GetUpdates fetches new messages from Telegram using long polling
func (b *Bot) GetUpdates(offset int) ([]Update, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.GetUpdatesWithContext(ctx, offset)
}

// GetUpdatesWithContext fetches new messages from Telegram using long polling with context support
func (b *Bot) GetUpdatesWithContext(ctx context.Context, offset int) ([]Update, error) {
	endpoint := fmt.Sprintf("/getUpdates?timeout=30&offset=%d", offset)

	resp, err := b.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Parse the result into []Update
	var updates []Update
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := json.Unmarshal(resultBytes, &updates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates: %w", err)
	}

	return updates, nil
}

// EditMessageTextOptions represents options for editing message text
type EditMessageTextOptions struct {
	ChatID                int64
	MessageID             int
	InlineMessageID       string
	Text                  string
	ParseMode             string
	Entities              []MessageEntity
	DisableWebPagePreview bool
	ReplyMarkup           interface{}
}

// EditMessageText edits the text of a message
func (b *Bot) EditMessageText(options *EditMessageTextOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.EditMessageTextWithContext(ctx, options)
}

// EditMessageTextWithContext edits the text of a message with context support
func (b *Bot) EditMessageTextWithContext(ctx context.Context, options *EditMessageTextOptions) error {
	reqBody := map[string]interface{}{
		"text": options.Text,
	}

	if options.ChatID != 0 {
		reqBody["chat_id"] = options.ChatID
	}
	if options.MessageID != 0 {
		reqBody["message_id"] = options.MessageID
	}
	if options.InlineMessageID != "" {
		reqBody["inline_message_id"] = options.InlineMessageID
	}
	if options.ParseMode != "" {
		reqBody["parse_mode"] = options.ParseMode
	}
	if len(options.Entities) > 0 {
		reqBody["entities"] = options.Entities
	}
	if options.DisableWebPagePreview {
		reqBody["disable_web_page_preview"] = true
	}
	if options.ReplyMarkup != nil {
		reqBody["reply_markup"] = options.ReplyMarkup
	}

	_, err := b.makeRequest(ctx, "POST", "/editMessageText", reqBody)
	return err
}

// DeleteMessageOptions represents options for deleting a message
type DeleteMessageOptions struct {
	ChatID    int64
	MessageID int
}

// DeleteMessage deletes a message
func (b *Bot) DeleteMessage(chatID int64, messageID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.DeleteMessageWithContext(ctx, chatID, messageID)
}

// DeleteMessageWithContext deletes a message with context support
func (b *Bot) DeleteMessageWithContext(ctx context.Context, chatID int64, messageID int) error {
	reqBody := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}

	_, err := b.makeRequest(ctx, "POST", "/deleteMessage", reqBody)
	return err
}

// SendPhotoOptions represents options for sending a photo
type SendPhotoOptions struct {
	ChatID                   int64
	Photo                    string // file_id, URL, or file path
	Caption                  string
	ParseMode                string
	CaptionEntities          []MessageEntity
	HasSpoiler               bool
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

// SendPhoto sends a photo
func (b *Bot) SendPhoto(options *SendPhotoOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendPhotoWithContext(ctx, options)
}

// SendPhotoWithContext sends a photo with context support
func (b *Bot) SendPhotoWithContext(ctx context.Context, options *SendPhotoOptions) error {
	reqBody := map[string]interface{}{
		"chat_id": options.ChatID,
		"photo":   options.Photo,
	}

	if options.Caption != "" {
		reqBody["caption"] = options.Caption
	}
	if options.ParseMode != "" {
		reqBody["parse_mode"] = options.ParseMode
	}
	if len(options.CaptionEntities) > 0 {
		reqBody["caption_entities"] = options.CaptionEntities
	}
	if options.HasSpoiler {
		reqBody["has_spoiler"] = true
	}
	if options.DisableNotification {
		reqBody["disable_notification"] = true
	}
	if options.ProtectContent {
		reqBody["protect_content"] = true
	}
	if options.ReplyToMessageID != 0 {
		reqBody["reply_to_message_id"] = options.ReplyToMessageID
	}
	if options.AllowSendingWithoutReply {
		reqBody["allow_sending_without_reply"] = true
	}
	if options.ReplyMarkup != nil {
		reqBody["reply_markup"] = options.ReplyMarkup
	}

	_, err := b.makeRequest(ctx, "POST", "/sendPhoto", reqBody)
	return err
}

// AnswerCallbackQuery answers a callback query
func (b *Bot) AnswerCallbackQuery(options *AnswerCallbackQueryOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.AnswerCallbackQueryWithContext(ctx, options)
}

// AnswerCallbackQueryWithContext answers a callback query with context support
func (b *Bot) AnswerCallbackQueryWithContext(ctx context.Context, options *AnswerCallbackQueryOptions) error {
	reqBody := map[string]interface{}{
		"callback_query_id": options.CallbackQueryID,
	}

	if options.Text != "" {
		reqBody["text"] = options.Text
	}
	if options.ShowAlert {
		reqBody["show_alert"] = true
	}
	if options.URL != "" {
		reqBody["url"] = options.URL
	}
	if options.CacheTime > 0 {
		reqBody["cache_time"] = options.CacheTime
	}

	_, err := b.makeRequest(ctx, "POST", "/answerCallbackQuery", reqBody)
	return err
}

// AnswerInlineQuery answers an inline query
func (b *Bot) AnswerInlineQuery(options *AnswerInlineQueryOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.AnswerInlineQueryWithContext(ctx, options)
}

// AnswerInlineQueryWithContext answers an inline query with context support
func (b *Bot) AnswerInlineQueryWithContext(ctx context.Context, options *AnswerInlineQueryOptions) error {
	reqBody := map[string]interface{}{
		"inline_query_id": options.InlineQueryID,
		"results":         options.Results,
	}

	if options.CacheTime > 0 {
		reqBody["cache_time"] = options.CacheTime
	}
	if options.IsPersonal {
		reqBody["is_personal"] = true
	}
	if options.NextOffset != "" {
		reqBody["next_offset"] = options.NextOffset
	}
	if options.SwitchPmText != "" {
		reqBody["switch_pm_text"] = options.SwitchPmText
	}
	if options.SwitchPmParameter != "" {
		reqBody["switch_pm_parameter"] = options.SwitchPmParameter
	}

	_, err := b.makeRequest(ctx, "POST", "/answerInlineQuery", reqBody)
	return err
}
