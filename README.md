# YouTube Downloader

A modern, elegant YouTube downloader built with Go and featuring a beautiful UI inspired by Apple and Stripe design aesthetics.

## Features

- 🎨 **Beautiful UI**: Clean, minimal design with smooth gradients and animations
- ⚡ **Fast Downloads**: Powered by yt-dlp and FFmpeg for reliable, high-speed downloads
- 📊 **Real-time Progress**: Live download progress updates via WebSocket
- 📁 **Auto File Explorer**: Automatically opens the download folder after completion
- 📜 **Download History**: Track all your downloaded videos
- 🎯 **Quality Selection**: Choose from multiple video quality options
- 🚀 **Large File Support**: Handle even the largest YouTube videos

## Prerequisites

1. **Go** (1.21 or higher)
2. **yt-dlp** - Install with:
   ```bash
   # Windows (with Python/pip)
   pip install yt-dlp
   
   # Or download directly from:
   # https://github.com/yt-dlp/yt-dlp/releases
   ```
3. **FFmpeg** - Already downloaded as mentioned

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd youtube-downloader
   ```
3. Download dependencies:
   ```bash
   go mod download
   ```

## Running the Application

1. Start the server:
   ```bash
   go run cmd/main.go
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## Usage

1. **Download a Video**:
   - Paste a YouTube URL in the input field
   - Select your preferred quality
   - Click "Download"
   - Watch the real-time progress
   - File Explorer opens automatically when complete

2. **View History**:
   - Scroll down to see all previously downloaded videos
   - Each entry shows the title, date, and status

## Project Structure

```
youtube-downloader/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── downloader/
│   │   └── manager.go       # Download logic and yt-dlp integration
│   ├── handlers/
│   │   ├── handlers.go      # HTTP request handlers
│   │   └── routes.go        # Route definitions
│   └── models/
│       └── models.go        # Data structures
├── static/
│   ├── css/
│   │   └── styles.css       # Modern CSS with gradients
│   └── js/
│       └── app.js           # Frontend JavaScript
├── templates/
│   └── index.html           # Main HTML template
├── downloads/               # Downloaded videos (created automatically)
├── go.mod
└── README.md
```

## Troubleshooting

### yt-dlp not found
Make sure yt-dlp is installed and available in your system PATH. You can verify with:
```bash
yt-dlp --version
```

### FFmpeg not found
Ensure FFmpeg is in your system PATH or update the FFmpeg location in `manager.go`.

### Port already in use
If port 8080 is already in use, you can change it in `cmd/main.go`.

## Building for Production

To create a standalone executable:
```bash
go build -o youtube-downloader.exe cmd/main.go
```

Then run:
```bash
./youtube-downloader.exe
```

## License

This project is for personal use only. Please respect YouTube's Terms of Service and copyright laws when downloading videos.