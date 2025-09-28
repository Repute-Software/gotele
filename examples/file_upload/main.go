package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gotele "github.com/repute-software/gotele/telegram"
)

func main() {
	// Load environment variables
	token := os.Getenv("BOT_DEV_TOKEN")
	if token == "" {
		log.Fatal("BOT_DEV_TOKEN is not set")
	}

	bot := gotele.NewBot(token)
	chatID := int64(123456789) // Replace with actual chat ID

	fmt.Println("=== File Upload & Media Handling Examples ===")

	// Example 1: Send document from file path
	fmt.Println("\n1. Sending document from file path...")
	documentOptions := &gotele.SendDocumentOptions{
		ChatID:    chatID,
		Document:  gotele.InputFile{FilePath: "example.txt"},
		Caption:   "Here's a text document! 📄",
		ParseMode: "Markdown",
	}

	err := bot.SendDocument(documentOptions)
	if err != nil {
		fmt.Printf("Error sending document: %v\n", err)
	}

	// Example 2: Send document from URL
	fmt.Println("\n2. Sending document from URL...")
	documentOptions2 := &gotele.SendDocumentOptions{
		ChatID:   chatID,
		Document: gotele.InputFile{URL: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"},
		Caption:  "PDF document from URL 📄",
	}

	err = bot.SendDocument(documentOptions2)
	if err != nil {
		fmt.Printf("Error sending document from URL: %v\n", err)
	}

	// Example 3: Send document from memory data
	fmt.Println("\n3. Sending document from memory data...")
	fileData := []byte("This is a test file content created in memory!")
	documentOptions3 := &gotele.SendDocumentOptions{
		ChatID: chatID,
		Document: gotele.InputFile{
			Data:     fileData,
			FileName: "memory_file.txt",
		},
		Caption: "File created from memory data! 💾",
	}

	err = bot.SendDocument(documentOptions3)
	if err != nil {
		fmt.Printf("Error sending document from memory: %v\n", err)
	}

	// Example 4: Send video with thumbnail
	fmt.Println("\n4. Sending video with thumbnail...")
	videoOptions := &gotele.SendVideoOptions{
		ChatID:            chatID,
		Video:             gotele.InputFile{URL: "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4"},
		Thumbnail:         gotele.InputFile{URL: "https://picsum.photos/320/240"},
		Caption:           "Sample video with custom thumbnail! 🎥",
		Duration:          60,
		Width:             1280,
		Height:            720,
		SupportsStreaming: true,
	}

	err = bot.SendVideo(videoOptions)
	if err != nil {
		fmt.Printf("Error sending video: %v\n", err)
	}

	// Example 5: Send audio file
	fmt.Println("\n5. Sending audio file...")
	audioOptions := &gotele.SendAudioOptions{
		ChatID:    chatID,
		Audio:     gotele.InputFile{URL: "https://www.soundjay.com/misc/sounds/bell-ringing-05.wav"},
		Caption:   "Audio file! 🎵",
		Duration:  5,
		Performer: "Test Artist",
		Title:     "Test Audio",
	}

	err = bot.SendAudio(audioOptions)
	if err != nil {
		fmt.Printf("Error sending audio: %v\n", err)
	}

	// Example 6: Send media group (album)
	fmt.Println("\n6. Sending media group (album)...")
	mediaGroup := []gotele.InputMedia{
		{
			Type:    "photo",
			Media:   "https://picsum.photos/400/300?random=1",
			Caption: "First photo in album 📸",
		},
		{
			Type:    "photo",
			Media:   "https://picsum.photos/400/300?random=2",
			Caption: "Second photo in album 📸",
		},
		{
			Type:    "photo",
			Media:   "https://picsum.photos/400/300?random=3",
			Caption: "Third photo in album 📸",
		},
	}

	mediaGroupOptions := &gotele.SendMediaGroupOptions{
		ChatID:              chatID,
		Media:               mediaGroup,
		DisableNotification: false,
		ProtectContent:      false,
	}

	err = bot.SendMediaGroup(mediaGroupOptions)
	if err != nil {
		fmt.Printf("Error sending media group: %v\n", err)
	}

	// Example 7: File download
	fmt.Println("\n7. Downloading file...")
	// First, we need a file ID from a previous message
	fileID := "BAADBAADrwADBREAAYag" // Replace with actual file ID
	file, err := bot.GetFile(fileID)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
	} else {
		fmt.Printf("File info: ID=%s, Size=%d bytes\n", file.FileID, file.FileSize)

		// Download file to local path
		err = bot.DownloadFileToPath(file, "./downloaded_file.jpg")
		if err != nil {
			fmt.Printf("Error downloading file: %v\n", err)
		} else {
			fmt.Println("File downloaded successfully!")
		}
	}

	// Example 8: File validation
	fmt.Println("\n8. File validation...")
	testFile := gotele.InputFile{FilePath: "example.txt"}

	// Get file size
	size, err := gotele.GetFileSize(testFile)
	if err != nil {
		fmt.Printf("Error getting file size: %v\n", err)
	} else {
		fmt.Printf("File size: %d bytes\n", size)

		// Validate file size
		err = gotele.ValidateFileSize(size, "document")
		if err != nil {
			fmt.Printf("File validation failed: %v\n", err)
		} else {
			fmt.Println("File size validation passed!")
		}
	}

	// Example 9: Send document with context and timeout
	fmt.Println("\n9. Sending document with context timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	documentOptions4 := &gotele.SendDocumentOptions{
		ChatID:              chatID,
		Document:            gotele.InputFile{URL: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"},
		Caption:             "Document with context timeout! ⏰",
		DisableNotification: true,
		ProtectContent:      true,
	}

	err = bot.SendDocumentWithContext(ctx, documentOptions4)
	if err != nil {
		fmt.Printf("Error sending document with context: %v\n", err)
	}

	// Example 10: Send document with reply markup
	fmt.Println("\n10. Sending document with inline keyboard...")
	keyboard := gotele.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotele.InlineKeyboardButton{
			{
				{Text: "Download", URL: "https://example.com/download"},
				{Text: "View Online", URL: "https://example.com/view"},
			},
		},
	}

	documentOptions5 := &gotele.SendDocumentOptions{
		ChatID:      chatID,
		Document:    gotele.InputFile{URL: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"},
		Caption:     "Document with inline keyboard! ⌨️",
		ReplyMarkup: keyboard,
	}

	err = bot.SendDocument(documentOptions5)
	if err != nil {
		fmt.Printf("Error sending document with keyboard: %v\n", err)
	}

	// Example 11: Send video with spoiler
	fmt.Println("\n11. Sending video with spoiler...")
	videoOptions2 := &gotele.SendVideoOptions{
		ChatID:     chatID,
		Video:      gotele.InputFile{URL: "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4"},
		Caption:    "Spoiler video! 🎬",
		HasSpoiler: true,
	}

	err = bot.SendVideo(videoOptions2)
	if err != nil {
		fmt.Printf("Error sending spoiler video: %v\n", err)
	}

	// Example 12: Send audio with thumbnail
	fmt.Println("\n12. Sending audio with thumbnail...")
	audioOptions2 := &gotele.SendAudioOptions{
		ChatID:    chatID,
		Audio:     gotele.InputFile{URL: "https://www.soundjay.com/misc/sounds/bell-ringing-05.wav"},
		Thumbnail: gotele.InputFile{URL: "https://picsum.photos/320/320"},
		Caption:   "Audio with thumbnail! 🎵🖼️",
		Duration:  5,
		Performer: "Test Artist",
		Title:     "Test Audio with Thumbnail",
	}

	err = bot.SendAudio(audioOptions2)
	if err != nil {
		fmt.Printf("Error sending audio with thumbnail: %v\n", err)
	}

	fmt.Println("\n=== All file upload examples completed ===")
}
