package gotele

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// makeMultipartRequest makes a multipart/form-data request for file uploads
func (b *Bot) makeMultipartRequest(ctx context.Context, endpoint string, fields map[string]string, files []FileUpload) (*APIResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	// Add files
	for _, file := range files {
		part, err := writer.CreateFormFile(file.FieldName, file.FileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file %s: %w", file.FieldName, err)
		}
		if _, err := part.Write(file.Data); err != nil {
			return nil, fmt.Errorf("failed to write file data for %s: %w", file.FieldName, err)
		}
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", b.BaseURL+endpoint, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make request
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

// prepareFileUpload prepares a file for upload
func (b *Bot) prepareFileUpload(inputFile InputFile, fieldName string) (FileUpload, error) {
	var data []byte
	var fileName string
	var mimeType string

	// Determine file source and read data
	if inputFile.FileID != "" {
		// File already uploaded to Telegram
		return FileUpload{
			FieldName: fieldName,
			FileName:  inputFile.FileName,
			Data:      []byte(inputFile.FileID),
			MimeType:  "text/plain",
		}, nil
	} else if inputFile.URL != "" {
		// File accessible via URL
		return FileUpload{
			FieldName: fieldName,
			FileName:  inputFile.FileName,
			Data:      []byte(inputFile.URL),
			MimeType:  "text/plain",
		}, nil
	} else if inputFile.FilePath != "" {
		// Local file
		fileData, err := os.ReadFile(inputFile.FilePath)
		if err != nil {
			return FileUpload{}, fmt.Errorf("failed to read file %s: %w", inputFile.FilePath, err)
		}
		data = fileData
		fileName = filepath.Base(inputFile.FilePath)
		if inputFile.FileName != "" {
			fileName = inputFile.FileName
		}
		mimeType = getMimeType(fileName)
	} else if len(inputFile.Data) > 0 {
		// File data in memory
		data = inputFile.Data
		fileName = inputFile.FileName
		if fileName == "" {
			fileName = "file"
		}
		mimeType = getMimeType(fileName)
	} else {
		return FileUpload{}, fmt.Errorf("no file data provided")
	}

	return FileUpload{
		FieldName: fieldName,
		FileName:  fileName,
		Data:      data,
		MimeType:  mimeType,
	}, nil
}

// getMimeType returns the MIME type based on file extension
func getMimeType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/avi"
	case ".mov":
		return "video/quicktime"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".txt":
		return "text/plain"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// SendDocumentOptions represents options for sending a document
type SendDocumentOptions struct {
	ChatID                      int64
	Document                    InputFile
	Thumbnail                   InputFile
	Caption                     string
	ParseMode                   string
	CaptionEntities             []MessageEntity
	DisableContentTypeDetection bool
	DisableNotification         bool
	ProtectContent              bool
	ReplyToMessageID            int
	AllowSendingWithoutReply    bool
	ReplyMarkup                 interface{}
}

// SendDocument sends a document
func (b *Bot) SendDocument(options *SendDocumentOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendDocumentWithContext(ctx, options)
}

// SendDocumentWithContext sends a document with context support
func (b *Bot) SendDocumentWithContext(ctx context.Context, options *SendDocumentOptions) error {
	// Prepare document file
	documentUpload, err := b.prepareFileUpload(options.Document, "document")
	if err != nil {
		return fmt.Errorf("failed to prepare document: %w", err)
	}

	files := []FileUpload{documentUpload}

	// Prepare thumbnail if provided
	if options.Thumbnail.FileID != "" || options.Thumbnail.URL != "" || options.Thumbnail.FilePath != "" || len(options.Thumbnail.Data) > 0 {
		thumbnailUpload, err := b.prepareFileUpload(options.Thumbnail, "thumbnail")
		if err != nil {
			return fmt.Errorf("failed to prepare thumbnail: %w", err)
		}
		files = append(files, thumbnailUpload)
	}

	// Prepare form fields
	fields := map[string]string{
		"chat_id": fmt.Sprintf("%d", options.ChatID),
	}

	if options.Caption != "" {
		fields["caption"] = options.Caption
	}
	if options.ParseMode != "" {
		fields["parse_mode"] = options.ParseMode
	}
	if options.DisableContentTypeDetection {
		fields["disable_content_type_detection"] = "true"
	}
	if options.DisableNotification {
		fields["disable_notification"] = "true"
	}
	if options.ProtectContent {
		fields["protect_content"] = "true"
	}
	if options.ReplyToMessageID != 0 {
		fields["reply_to_message_id"] = fmt.Sprintf("%d", options.ReplyToMessageID)
	}
	if options.AllowSendingWithoutReply {
		fields["allow_sending_without_reply"] = "true"
	}
	if options.ReplyMarkup != nil {
		// Convert reply markup to JSON
		markupJSON, err := json.Marshal(options.ReplyMarkup)
		if err != nil {
			return fmt.Errorf("failed to marshal reply markup: %w", err)
		}
		fields["reply_markup"] = string(markupJSON)
	}

	_, err = b.makeMultipartRequest(ctx, "/sendDocument", fields, files)
	return err
}

// SendVideoOptions represents options for sending a video
type SendVideoOptions struct {
	ChatID                   int64
	Video                    InputFile
	Duration                 int
	Width                    int
	Height                   int
	Thumbnail                InputFile
	Caption                  string
	ParseMode                string
	CaptionEntities          []MessageEntity
	HasSpoiler               bool
	SupportsStreaming        bool
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

// SendVideo sends a video
func (b *Bot) SendVideo(options *SendVideoOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendVideoWithContext(ctx, options)
}

// SendVideoWithContext sends a video with context support
func (b *Bot) SendVideoWithContext(ctx context.Context, options *SendVideoOptions) error {
	// Prepare video file
	videoUpload, err := b.prepareFileUpload(options.Video, "video")
	if err != nil {
		return fmt.Errorf("failed to prepare video: %w", err)
	}

	files := []FileUpload{videoUpload}

	// Prepare thumbnail if provided
	if options.Thumbnail.FileID != "" || options.Thumbnail.URL != "" || options.Thumbnail.FilePath != "" || len(options.Thumbnail.Data) > 0 {
		thumbnailUpload, err := b.prepareFileUpload(options.Thumbnail, "thumbnail")
		if err != nil {
			return fmt.Errorf("failed to prepare thumbnail: %w", err)
		}
		files = append(files, thumbnailUpload)
	}

	// Prepare form fields
	fields := map[string]string{
		"chat_id": fmt.Sprintf("%d", options.ChatID),
	}

	if options.Duration != 0 {
		fields["duration"] = fmt.Sprintf("%d", options.Duration)
	}
	if options.Width != 0 {
		fields["width"] = fmt.Sprintf("%d", options.Width)
	}
	if options.Height != 0 {
		fields["height"] = fmt.Sprintf("%d", options.Height)
	}
	if options.Caption != "" {
		fields["caption"] = options.Caption
	}
	if options.ParseMode != "" {
		fields["parse_mode"] = options.ParseMode
	}
	if options.HasSpoiler {
		fields["has_spoiler"] = "true"
	}
	if options.SupportsStreaming {
		fields["supports_streaming"] = "true"
	}
	if options.DisableNotification {
		fields["disable_notification"] = "true"
	}
	if options.ProtectContent {
		fields["protect_content"] = "true"
	}
	if options.ReplyToMessageID != 0 {
		fields["reply_to_message_id"] = fmt.Sprintf("%d", options.ReplyToMessageID)
	}
	if options.AllowSendingWithoutReply {
		fields["allow_sending_without_reply"] = "true"
	}
	if options.ReplyMarkup != nil {
		markupJSON, err := json.Marshal(options.ReplyMarkup)
		if err != nil {
			return fmt.Errorf("failed to marshal reply markup: %w", err)
		}
		fields["reply_markup"] = string(markupJSON)
	}

	_, err = b.makeMultipartRequest(ctx, "/sendVideo", fields, files)
	return err
}

// SendAudioOptions represents options for sending an audio file
type SendAudioOptions struct {
	ChatID                   int64
	Audio                    InputFile
	Caption                  string
	ParseMode                string
	CaptionEntities          []MessageEntity
	Duration                 int
	Performer                string
	Title                    string
	Thumbnail                InputFile
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

// SendAudio sends an audio file
func (b *Bot) SendAudio(options *SendAudioOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendAudioWithContext(ctx, options)
}

// SendAudioWithContext sends an audio file with context support
func (b *Bot) SendAudioWithContext(ctx context.Context, options *SendAudioOptions) error {
	// Prepare audio file
	audioUpload, err := b.prepareFileUpload(options.Audio, "audio")
	if err != nil {
		return fmt.Errorf("failed to prepare audio: %w", err)
	}

	files := []FileUpload{audioUpload}

	// Prepare thumbnail if provided
	if options.Thumbnail.FileID != "" || options.Thumbnail.URL != "" || options.Thumbnail.FilePath != "" || len(options.Thumbnail.Data) > 0 {
		thumbnailUpload, err := b.prepareFileUpload(options.Thumbnail, "thumbnail")
		if err != nil {
			return fmt.Errorf("failed to prepare thumbnail: %w", err)
		}
		files = append(files, thumbnailUpload)
	}

	// Prepare form fields
	fields := map[string]string{
		"chat_id": fmt.Sprintf("%d", options.ChatID),
	}

	if options.Caption != "" {
		fields["caption"] = options.Caption
	}
	if options.ParseMode != "" {
		fields["parse_mode"] = options.ParseMode
	}
	if options.Duration != 0 {
		fields["duration"] = fmt.Sprintf("%d", options.Duration)
	}
	if options.Performer != "" {
		fields["performer"] = options.Performer
	}
	if options.Title != "" {
		fields["title"] = options.Title
	}
	if options.DisableNotification {
		fields["disable_notification"] = "true"
	}
	if options.ProtectContent {
		fields["protect_content"] = "true"
	}
	if options.ReplyToMessageID != 0 {
		fields["reply_to_message_id"] = fmt.Sprintf("%d", options.ReplyToMessageID)
	}
	if options.AllowSendingWithoutReply {
		fields["allow_sending_without_reply"] = "true"
	}
	if options.ReplyMarkup != nil {
		markupJSON, err := json.Marshal(options.ReplyMarkup)
		if err != nil {
			return fmt.Errorf("failed to marshal reply markup: %w", err)
		}
		fields["reply_markup"] = string(markupJSON)
	}

	_, err = b.makeMultipartRequest(ctx, "/sendAudio", fields, files)
	return err
}

// GetFileOptions represents options for getting file information
type GetFileOptions struct {
	FileID string
}

// GetFile gets file information
func (b *Bot) GetFile(fileID string) (*File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.GetFileWithContext(ctx, fileID)
}

// GetFileWithContext gets file information with context support
func (b *Bot) GetFileWithContext(ctx context.Context, fileID string) (*File, error) {
	endpoint := fmt.Sprintf("/getFile?file_id=%s", fileID)

	resp, err := b.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Parse the result into File
	var file File
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := json.Unmarshal(resultBytes, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %w", err)
	}

	return &file, nil
}

// DownloadFile downloads a file from Telegram servers
func (b *Bot) DownloadFile(file *File) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.DownloadFileWithContext(ctx, file)
}

// DownloadFileWithContext downloads a file from Telegram servers with context support
func (b *Bot) DownloadFileWithContext(ctx context.Context, file *File) ([]byte, error) {
	if file.FilePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	// Construct download URL
	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.Token, file.FilePath)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	// Make request
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Read file data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return data, nil
}

// DownloadFileToPath downloads a file and saves it to a local path
func (b *Bot) DownloadFileToPath(file *File, localPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.DownloadFileToPathWithContext(ctx, file, localPath)
}

// DownloadFileToPathWithContext downloads a file and saves it to a local path with context support
func (b *Bot) DownloadFileToPathWithContext(ctx context.Context, file *File, localPath string) error {
	data, err := b.DownloadFileWithContext(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file to %s: %w", localPath, err)
	}

	return nil
}

// SendMediaGroupOptions represents options for sending a media group
type SendMediaGroupOptions struct {
	ChatID                   int64
	Media                    []InputMedia
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
}

// SendMediaGroup sends a group of media files as an album
func (b *Bot) SendMediaGroup(options *SendMediaGroupOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.SendMediaGroupWithContext(ctx, options)
}

// SendMediaGroupWithContext sends a group of media files as an album with context support
func (b *Bot) SendMediaGroupWithContext(ctx context.Context, options *SendMediaGroupOptions) error {
	// Convert InputMedia to JSON
	mediaJSON, err := json.Marshal(options.Media)
	if err != nil {
		return fmt.Errorf("failed to marshal media: %w", err)
	}

	reqBody := map[string]interface{}{
		"chat_id": options.ChatID,
		"media":   string(mediaJSON),
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

	_, err = b.makeRequest(ctx, "POST", "/sendMediaGroup", reqBody)
	return err
}

// File validation constants
const (
	MaxFileSize     = 50 * 1024 * 1024       // 50MB
	MaxPhotoSize    = 10 * 1024 * 1024       // 10MB
	MaxVideoSize    = 50 * 1024 * 1024       // 50MB
	MaxAudioSize    = 50 * 1024 * 1024       // 50MB
	MaxDocumentSize = 2 * 1024 * 1024 * 1024 // 2GB
)

// ValidateFileSize validates file size against Telegram limits
func ValidateFileSize(fileSize int64, fileType string) error {
	var maxSize int64

	switch fileType {
	case "photo":
		maxSize = MaxPhotoSize
	case "video":
		maxSize = MaxVideoSize
	case "audio":
		maxSize = MaxAudioSize
	case "document":
		maxSize = MaxDocumentSize
	default:
		maxSize = MaxFileSize
	}

	if fileSize > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes for %s", fileSize, maxSize, fileType)
	}

	return nil
}

// GetFileSize returns the size of a file
func GetFileSize(inputFile InputFile) (int64, error) {
	if inputFile.FilePath != "" {
		info, err := os.Stat(inputFile.FilePath)
		if err != nil {
			return 0, fmt.Errorf("failed to get file info: %w", err)
		}
		return info.Size(), nil
	} else if len(inputFile.Data) > 0 {
		return int64(len(inputFile.Data)), nil
	}

	return 0, fmt.Errorf("cannot determine file size")
}
