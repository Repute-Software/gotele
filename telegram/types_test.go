package gotele

import (
	"testing"
)

func TestUser(t *testing.T) {
	user := User{
		ID:        123456789,
		IsBot:     false,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
	}

	if user.ID != 123456789 {
		t.Errorf("Expected ID 123456789, got %d", user.ID)
	}
	if user.IsBot {
		t.Error("Expected IsBot to be false")
	}
	if user.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got %s", user.FirstName)
	}
}

func TestChat(t *testing.T) {
	chat := Chat{
		ID:    -1001234567890,
		Type:  "supergroup",
		Title: "Test Group",
	}

	if chat.ID != -1001234567890 {
		t.Errorf("Expected ID -1001234567890, got %d", chat.ID)
	}
	if chat.Type != "supergroup" {
		t.Errorf("Expected Type 'supergroup', got %s", chat.Type)
	}
	if chat.Title != "Test Group" {
		t.Errorf("Expected Title 'Test Group', got %s", chat.Title)
	}
}

func TestMessage(t *testing.T) {
	user := User{
		ID:        123456789,
		IsBot:     false,
		FirstName: "John",
	}

	chat := Chat{
		ID:   123456789,
		Type: "private",
	}

	message := Message{
		MessageID: 1,
		From:      &user,
		Chat:      chat,
		Text:      "Hello, World!",
		Date:      1640995200,
	}

	if message.MessageID != 1 {
		t.Errorf("Expected MessageID 1, got %d", message.MessageID)
	}
	if message.From == nil {
		t.Error("Expected From to be set")
	}
	if message.Text != "Hello, World!" {
		t.Errorf("Expected Text 'Hello, World!', got %s", message.Text)
	}
}

func TestMessageEntity(t *testing.T) {
	entity := MessageEntity{
		Type:   "bold",
		Offset: 0,
		Length: 5,
	}

	if entity.Type != "bold" {
		t.Errorf("Expected Type 'bold', got %s", entity.Type)
	}
	if entity.Offset != 0 {
		t.Errorf("Expected Offset 0, got %d", entity.Offset)
	}
	if entity.Length != 5 {
		t.Errorf("Expected Length 5, got %d", entity.Length)
	}
}

func TestPhotoSize(t *testing.T) {
	photo := PhotoSize{
		FileID:       "BAADBAADrwADBREAAYag",
		FileUniqueID: "AgADrwADBREAAYag",
		Width:        1280,
		Height:       720,
		FileSize:     12345,
	}

	if photo.FileID != "BAADBAADrwADBREAAYag" {
		t.Errorf("Expected FileID 'BAADBAADrwADBREAAYag', got %s", photo.FileID)
	}
	if photo.Width != 1280 {
		t.Errorf("Expected Width 1280, got %d", photo.Width)
	}
	if photo.Height != 720 {
		t.Errorf("Expected Height 720, got %d", photo.Height)
	}
}

func TestInlineKeyboardMarkup(t *testing.T) {
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "btn1"},
				{Text: "Button 2", CallbackData: "btn2"},
			},
			{
				{Text: "Button 3", URL: "https://example.com"},
			},
		},
	}

	if len(keyboard.InlineKeyboard) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(keyboard.InlineKeyboard))
	}
	if len(keyboard.InlineKeyboard[0]) != 2 {
		t.Errorf("Expected 2 buttons in first row, got %d", len(keyboard.InlineKeyboard[0]))
	}
	if keyboard.InlineKeyboard[0][0].Text != "Button 1" {
		t.Errorf("Expected first button text 'Button 1', got %s", keyboard.InlineKeyboard[0][0].Text)
	}
}

func TestReplyKeyboardMarkup(t *testing.T) {
	keyboard := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Button 1"},
				{Text: "Button 2"},
			},
			{
				{Text: "Button 3", RequestContact: true},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	if len(keyboard.Keyboard) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(keyboard.Keyboard))
	}
	if !keyboard.ResizeKeyboard {
		t.Error("Expected ResizeKeyboard to be true")
	}
	if !keyboard.OneTimeKeyboard {
		t.Error("Expected OneTimeKeyboard to be true")
	}
}

func TestLocation(t *testing.T) {
	location := Location{
		Longitude: 37.7749,
		Latitude:  -122.4194,
	}

	if location.Longitude != 37.7749 {
		t.Errorf("Expected Longitude 37.7749, got %f", location.Longitude)
	}
	if location.Latitude != -122.4194 {
		t.Errorf("Expected Latitude -122.4194, got %f", location.Latitude)
	}
}

func TestContact(t *testing.T) {
	contact := Contact{
		PhoneNumber: "+1234567890",
		FirstName:   "John",
		LastName:    "Doe",
		UserID:      123456789,
	}

	if contact.PhoneNumber != "+1234567890" {
		t.Errorf("Expected PhoneNumber '+1234567890', got %s", contact.PhoneNumber)
	}
	if contact.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got %s", contact.FirstName)
	}
	if contact.UserID != 123456789 {
		t.Errorf("Expected UserID 123456789, got %d", contact.UserID)
	}
}

func TestPoll(t *testing.T) {
	poll := Poll{
		ID:       "poll123",
		Question: "What is your favorite color?",
		Options: []PollOption{
			{Text: "Red", VoterCount: 5},
			{Text: "Blue", VoterCount: 3},
			{Text: "Green", VoterCount: 2},
		},
		TotalVoterCount: 10,
		IsClosed:        false,
		IsAnonymous:     true,
		Type:            "regular",
	}

	if poll.ID != "poll123" {
		t.Errorf("Expected ID 'poll123', got %s", poll.ID)
	}
	if poll.Question != "What is your favorite color?" {
		t.Errorf("Expected Question 'What is your favorite color?', got %s", poll.Question)
	}
	if len(poll.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(poll.Options))
	}
	if poll.TotalVoterCount != 10 {
		t.Errorf("Expected TotalVoterCount 10, got %d", poll.TotalVoterCount)
	}
}

func TestSticker(t *testing.T) {
	sticker := Sticker{
		FileID:       "CAADBAADrwADBREAAYag",
		FileUniqueID: "AgADrwADBREAAYag",
		Type:         "regular",
		Width:        512,
		Height:       512,
		IsAnimated:   false,
		IsVideo:      false,
		Emoji:        "😀",
	}

	if sticker.FileID != "CAADBAADrwADBREAAYag" {
		t.Errorf("Expected FileID 'CAADBAADrwADBREAAYag', got %s", sticker.FileID)
	}
	if sticker.Type != "regular" {
		t.Errorf("Expected Type 'regular', got %s", sticker.Type)
	}
	if sticker.Width != 512 {
		t.Errorf("Expected Width 512, got %d", sticker.Width)
	}
	if sticker.IsAnimated {
		t.Error("Expected IsAnimated to be false")
	}
	if sticker.Emoji != "😀" {
		t.Errorf("Expected Emoji '😀', got %s", sticker.Emoji)
	}
}
