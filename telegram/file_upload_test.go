package gotele

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestInputFile(t *testing.T) {
	// Test InputFile with file path
	inputFile := InputFile{
		FilePath: "test.txt",
		FileName: "custom_name.txt",
	}

	if inputFile.FilePath != "test.txt" {
		t.Errorf("Expected FilePath 'test.txt', got %s", inputFile.FilePath)
	}
	if inputFile.FileName != "custom_name.txt" {
		t.Errorf("Expected FileName 'custom_name.txt', got %s", inputFile.FileName)
	}

	// Test InputFile with data
	data := []byte("test data")
	inputFile2 := InputFile{
		Data:     data,
		FileName: "data.txt",
	}

	if len(inputFile2.Data) != len(data) {
		t.Errorf("Expected data length %d, got %d", len(data), len(inputFile2.Data))
	}
}

func TestFileUpload(t *testing.T) {
	upload := FileUpload{
		FieldName: "document",
		FileName:  "test.txt",
		Data:      []byte("test content"),
		MimeType:  "text/plain",
	}

	if upload.FieldName != "document" {
		t.Errorf("Expected FieldName 'document', got %s", upload.FieldName)
	}
	if upload.FileName != "test.txt" {
		t.Errorf("Expected FileName 'test.txt', got %s", upload.FileName)
	}
	if upload.MimeType != "text/plain" {
		t.Errorf("Expected MimeType 'text/plain', got %s", upload.MimeType)
	}
}

func TestInputMedia(t *testing.T) {
	media := InputMedia{
		Type:       "photo",
		Media:      "file_id_123",
		Caption:    "Test photo",
		ParseMode:  "Markdown",
		HasSpoiler: false,
	}

	if media.Type != "photo" {
		t.Errorf("Expected Type 'photo', got %s", media.Type)
	}
	if media.Media != "file_id_123" {
		t.Errorf("Expected Media 'file_id_123', got %s", media.Media)
	}
	if media.Caption != "Test photo" {
		t.Errorf("Expected Caption 'Test photo', got %s", media.Caption)
	}
}

func TestInputMediaPhoto(t *testing.T) {
	photo := InputMediaPhoto{
		Type:       "photo",
		Media:      "file_id_456",
		Caption:    "Beautiful sunset",
		ParseMode:  "HTML",
		HasSpoiler: true,
	}

	if photo.Type != "photo" {
		t.Errorf("Expected Type 'photo', got %s", photo.Type)
	}
	if photo.HasSpoiler != true {
		t.Error("Expected HasSpoiler to be true")
	}
}

func TestInputMediaVideo(t *testing.T) {
	video := InputMediaVideo{
		Type:              "video",
		Media:             "file_id_789",
		Width:             1920,
		Height:            1080,
		Duration:          120,
		SupportsStreaming: true,
		HasSpoiler:        false,
	}

	if video.Type != "video" {
		t.Errorf("Expected Type 'video', got %s", video.Type)
	}
	if video.Width != 1920 {
		t.Errorf("Expected Width 1920, got %d", video.Width)
	}
	if video.Height != 1080 {
		t.Errorf("Expected Height 1080, got %d", video.Height)
	}
	if video.Duration != 120 {
		t.Errorf("Expected Duration 120, got %d", video.Duration)
	}
	if !video.SupportsStreaming {
		t.Error("Expected SupportsStreaming to be true")
	}
}

func TestInputMediaDocument(t *testing.T) {
	document := InputMediaDocument{
		Type:                        "document",
		Media:                       "file_id_doc",
		Caption:                     "Important document",
		ParseMode:                   "Markdown",
		DisableContentTypeDetection: true,
	}

	if document.Type != "document" {
		t.Errorf("Expected Type 'document', got %s", document.Type)
	}
	if !document.DisableContentTypeDetection {
		t.Error("Expected DisableContentTypeDetection to be true")
	}
}

func TestInputMediaAudio(t *testing.T) {
	audio := InputMediaAudio{
		Type:      "audio",
		Media:     "file_id_audio",
		Duration:  180,
		Performer: "Test Artist",
		Title:     "Test Song",
	}

	if audio.Type != "audio" {
		t.Errorf("Expected Type 'audio', got %s", audio.Type)
	}
	if audio.Duration != 180 {
		t.Errorf("Expected Duration 180, got %d", audio.Duration)
	}
	if audio.Performer != "Test Artist" {
		t.Errorf("Expected Performer 'Test Artist', got %s", audio.Performer)
	}
	if audio.Title != "Test Song" {
		t.Errorf("Expected Title 'Test Song', got %s", audio.Title)
	}
}

func TestSendDocumentOptions(t *testing.T) {
	options := SendDocumentOptions{
		ChatID:                      123456789,
		Document:                    InputFile{FilePath: "test.pdf"},
		Caption:                     "Test document",
		ParseMode:                   "Markdown",
		DisableContentTypeDetection: true,
		DisableNotification:         true,
		ProtectContent:              true,
		ReplyToMessageID:            1,
		AllowSendingWithoutReply:    true,
	}

	if options.ChatID != 123456789 {
		t.Errorf("Expected ChatID 123456789, got %d", options.ChatID)
	}
	if options.Caption != "Test document" {
		t.Errorf("Expected Caption 'Test document', got %s", options.Caption)
	}
	if !options.DisableContentTypeDetection {
		t.Error("Expected DisableContentTypeDetection to be true")
	}
	if !options.DisableNotification {
		t.Error("Expected DisableNotification to be true")
	}
	if !options.ProtectContent {
		t.Error("Expected ProtectContent to be true")
	}
}

func TestSendVideoOptions(t *testing.T) {
	options := SendVideoOptions{
		ChatID:              123456789,
		Video:               InputFile{FilePath: "test.mp4"},
		Duration:            60,
		Width:               1280,
		Height:              720,
		Caption:             "Test video",
		HasSpoiler:          true,
		SupportsStreaming:   true,
		DisableNotification: false,
		ProtectContent:      false,
	}

	if options.ChatID != 123456789 {
		t.Errorf("Expected ChatID 123456789, got %d", options.ChatID)
	}
	if options.Duration != 60 {
		t.Errorf("Expected Duration 60, got %d", options.Duration)
	}
	if options.Width != 1280 {
		t.Errorf("Expected Width 1280, got %d", options.Width)
	}
	if options.Height != 720 {
		t.Errorf("Expected Height 720, got %d", options.Height)
	}
	if !options.HasSpoiler {
		t.Error("Expected HasSpoiler to be true")
	}
	if !options.SupportsStreaming {
		t.Error("Expected SupportsStreaming to be true")
	}
}

func TestSendAudioOptions(t *testing.T) {
	options := SendAudioOptions{
		ChatID:              123456789,
		Audio:               InputFile{FilePath: "test.mp3"},
		Caption:             "Test audio",
		Duration:            120,
		Performer:           "Test Artist",
		Title:               "Test Song",
		DisableNotification: true,
		ProtectContent:      true,
	}

	if options.ChatID != 123456789 {
		t.Errorf("Expected ChatID 123456789, got %d", options.ChatID)
	}
	if options.Duration != 120 {
		t.Errorf("Expected Duration 120, got %d", options.Duration)
	}
	if options.Performer != "Test Artist" {
		t.Errorf("Expected Performer 'Test Artist', got %s", options.Performer)
	}
	if options.Title != "Test Song" {
		t.Errorf("Expected Title 'Test Song', got %s", options.Title)
	}
}

func TestSendMediaGroupOptions(t *testing.T) {
	media := []InputMedia{
		{Type: "photo", Media: "file_id_1", Caption: "Photo 1"},
		{Type: "photo", Media: "file_id_2", Caption: "Photo 2"},
	}

	options := SendMediaGroupOptions{
		ChatID:                   123456789,
		Media:                    media,
		DisableNotification:      true,
		ProtectContent:           true,
		ReplyToMessageID:         1,
		AllowSendingWithoutReply: true,
	}

	if options.ChatID != 123456789 {
		t.Errorf("Expected ChatID 123456789, got %d", options.ChatID)
	}
	if len(options.Media) != 2 {
		t.Errorf("Expected 2 media items, got %d", len(options.Media))
	}
	if !options.DisableNotification {
		t.Error("Expected DisableNotification to be true")
	}
}

func TestValidateFileSize(t *testing.T) {
	// Test valid file sizes
	err := ValidateFileSize(5*1024*1024, "photo") // 5MB photo
	if err != nil {
		t.Errorf("Expected no error for 5MB photo, got %v", err)
	}

	err = ValidateFileSize(25*1024*1024, "video") // 25MB video
	if err != nil {
		t.Errorf("Expected no error for 25MB video, got %v", err)
	}

	err = ValidateFileSize(1*1024*1024*1024, "document") // 1GB document
	if err != nil {
		t.Errorf("Expected no error for 1GB document, got %v", err)
	}

	// Test invalid file sizes
	err = ValidateFileSize(15*1024*1024, "photo") // 15MB photo (exceeds 10MB limit)
	if err == nil {
		t.Error("Expected error for 15MB photo")
	}

	err = ValidateFileSize(100*1024*1024, "video") // 100MB video (exceeds 50MB limit)
	if err == nil {
		t.Error("Expected error for 100MB video")
	}
}

func TestGetFileSize(t *testing.T) {
	// Test with file path (create a temporary file)
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := []byte("test content")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	inputFile := InputFile{FilePath: tmpFile.Name()}
	size, err := GetFileSize(inputFile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if size != int64(len(testData)) {
		t.Errorf("Expected size %d, got %d", len(testData), size)
	}

	// Test with data
	inputFile2 := InputFile{Data: testData}
	size2, err := GetFileSize(inputFile2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if size2 != int64(len(testData)) {
		t.Errorf("Expected size %d, got %d", len(testData), size2)
	}

	// Test with no data
	inputFile3 := InputFile{}
	_, err = GetFileSize(inputFile3)
	if err == nil {
		t.Error("Expected error for empty input file")
	}
}

func TestFileUploadWithContext(t *testing.T) {
	bot := NewBot("test_token")

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	options := &SendDocumentOptions{
		ChatID:   123456789,
		Document: InputFile{FilePath: "nonexistent.txt"},
	}

	err := bot.SendDocumentWithContext(ctx, options)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}

func TestFileUploadWithTimeout(t *testing.T) {
	bot := NewBot("test_token")

	// Test with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(1 * time.Millisecond)

	options := &SendDocumentOptions{
		ChatID:   123456789,
		Document: InputFile{FilePath: "test.txt"},
	}

	err := bot.SendDocumentWithContext(ctx, options)
	if err == nil {
		t.Error("Expected error for timed out context")
	}
}
