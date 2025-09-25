package gotele

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewBot creates a new Bot instance
func NewBot(token string) *Bot {
	return &Bot{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Client:  &http.Client{},
	}
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// SendMessage sends a message to a specific chat
func (b *Bot) SendMessage(chatID int64, text string) error {
	url := b.BaseURL + "/sendMessage"

	reqBody := sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := b.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", resp.Status)
	}

	return nil
}

// GetUpdates fetches new messages from Telegram using long polling
func (b *Bot) GetUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?timeout=30&offset=%d", b.BaseURL, offset)

	resp, err := b.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram API error: %s", resp.Status)
	}

	var updatesResp GetUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatesResp); err != nil {
		return nil, err
	}

	if !updatesResp.Ok {
		return nil, fmt.Errorf("telegram API returned ok=false")
	}

	return updatesResp.Result, nil
}
