## File uploads and downloads

`InputFile` supports multiple sources:

- `FileID`: reuse a file already on Telegram servers
- `URL`: point to a remote file
- `FilePath`: read from local path
- `Data`: provide bytes in-memory (with optional `FileName`)

### Send a document

```go
opts := &gotele.SendDocumentOptions{
    ChatID:   chatID,
    Document: gotele.InputFile{FilePath: "example.pdf"},
    Caption:  "Here you go",
}
_ = bot.SendDocument(opts)
```

### Send a video with thumbnail

```go
opts := &gotele.SendVideoOptions{
    ChatID:    chatID,
    Video:     gotele.InputFile{URL: "https://.../video.mp4"},
    Thumbnail: gotele.InputFile{FilePath: "thumb.jpg"},
}
_ = bot.SendVideo(opts)
```

### Send audio

```go
_ = bot.SendAudio(&gotele.SendAudioOptions{ChatID: chatID, Audio: gotele.InputFile{Data: data, FileName: "track.mp3"}})
```

### Media groups

```go
media := []gotele.InputMedia{
    {Type: "photo", Media: "attach://photo1"},
    {Type: "photo", Media: "attach://photo2"},
}
_ = bot.SendMediaGroup(&gotele.SendMediaGroupOptions{ChatID: chatID, Media: media})
```

### File download

```go
file, _ := bot.GetFile(fileID)
data, _ := bot.DownloadFile(file)
_ = bot.DownloadFileToPath(file, "downloads/file.bin")
```

### Validation helpers

```go
size, _ := gotele.GetFileSize(gotele.InputFile{FilePath: "example.pdf"})
_ = gotele.ValidateFileSize(size, "document")
```

