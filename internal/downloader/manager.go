package downloader

import (
	"Go-Utilities/internal/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	downloads   map[string]*Download
	subscribers []chan models.ProgressUpdate
	mu          sync.RWMutex
}

type Download struct {
	ID       string
	URL      string
	Title    string
	Status   string
	Progress float64
	Speed    string
	ETA      string
}

func NewManager() *Manager {
	return &Manager{
		downloads:   make(map[string]*Download),
		subscribers: []chan models.ProgressUpdate{},
	}
}

func getYtDlpPath() (string, error) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Build absolute path to yt-dlp.exe
	ytDlpPath := filepath.Join(wd, "yt-dlp.exe")

	// Check if it exists
	if _, err := os.Stat(ytDlpPath); os.IsNotExist(err) {
		return "", fmt.Errorf("yt-dlp.exe not found at %s", ytDlpPath)
	}

	return ytDlpPath, nil
}

func (m *Manager) getFFmpegPath() (string, error) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Build absolute path to ffmpeg.exe
	ffmpegPath := filepath.Join(wd, "ffmpeg.exe")

	// Check if it exists
	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
		return "", fmt.Errorf("ffmpeg.exe not found at %s", ffmpegPath)
	}

	return ffmpegPath, nil
}

func (m *Manager) TestYtDlp() error {
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return err
	}

	// Test with version command
	cmd := exec.Command(ytDlpPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("yt-dlp test failed: %v", err)
	}

	version := strings.TrimSpace(string(output))
	log.Printf("yt-dlp version: %s", version)

	// Check if version looks old (very basic check)
	if len(version) > 0 && version < "2024" {
		log.Printf("WARNING: yt-dlp version may be outdated. Consider updating from https://github.com/yt-dlp/yt-dlp/releases")
	}

	return nil
}

func (m *Manager) findDownloadedFile(dir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			// Skip fragment files and other temp files
			name := strings.ToLower(filepath.Base(file))
			if !strings.Contains(name, ".f") && !strings.Contains(name, ".temp") {
				return file, nil
			}
		}
	}

	return "", fmt.Errorf("no video file found in %s", dir)
}

func (m *Manager) addResolutionToFilename(filePath, quality string) string {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Add resolution to filename if not already present
	if quality != "" && quality != "best" && !strings.Contains(nameWithoutExt, quality) {
		newFilename := fmt.Sprintf("%s [%s]%s", nameWithoutExt, quality, ext)
		return filepath.Join(dir, newFilename)
	}

	return filePath
}

func (m *Manager) getUniqueFilePath(dir, filename string) string {
	basePath := filepath.Join(dir, filename)

	// If file doesn't exist, return original path
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return basePath
	}

	// File exists, need to create unique name
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	counter := 1
	for {
		// Create new filename with suffix
		newFilename := fmt.Sprintf("%s-%d%s", nameWithoutExt, counter, ext)
		newPath := filepath.Join(dir, newFilename)

		// Check if this unique name is available
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			log.Printf("Duplicate file detected. Saving as: %s", newFilename)
			return newPath
		}

		counter++

		// Safety check to prevent infinite loop
		if counter > 1000 {
			// If we somehow hit 1000 duplicates, use timestamp
			timestamp := time.Now().Format("20060102-150405")
			newFilename := fmt.Sprintf("%s-%s%s", nameWithoutExt, timestamp, ext)
			return filepath.Join(dir, newFilename)
		}
	}
}

func (m *Manager) copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func (m *Manager) StartDownload(url, quality string) string {
	downloadID := fmt.Sprintf("dl_%d", time.Now().Unix())

	go m.download(downloadID, url, quality)

	return downloadID
}

func (m *Manager) StartMp3Convert(url string) string {
	downloadID := fmt.Sprintf("mp3_%d", time.Now().Unix())

	go m.convertToMp3(downloadID, url)

	return downloadID
}

func (m *Manager) download(id, url, quality string) {
	m.mu.Lock()
	m.downloads[id] = &Download{
		ID:     id,
		URL:    url,
		Status: "starting",
	}
	m.mu.Unlock()

	// Create downloads directory
	downloadDir := "downloads"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create download directory: %v", err))
		return
	}

	// Prepare yt-dlp command - use temp directory first
	tempDir := filepath.Join(os.TempDir(), "ytdownloader")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create temp directory: %v", err))
		return
	}

	outputPath := filepath.Join(tempDir, "%(title)s.%(ext)s")
	// Get FFmpeg path
	ffmpegPath, err := m.getFFmpegPath()
	if err != nil {
		log.Printf("Warning: FFmpeg not found, audio merging may not work: %v", err)
	}

	args := []string{
		"--no-warnings",
		"--newline",
		"--progress",
		"-o", outputPath,
		"--merge-output-format", "mp4",
		"--embed-metadata",
		"--no-playlist",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"--referer", "https://www.youtube.com/",
		"--add-header", "Accept-Language:en-US,en;q=0.9",
		"--add-header", "Accept-Encoding:gzip, deflate, br, zstd",
		"--add-header", "Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"--add-header", "Cache-Control:no-cache",
		"--add-header", "Pragma:no-cache",
		"--add-header", "Sec-Ch-Ua:\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
		"--add-header", "Sec-Ch-Ua-Mobile:?0",
		"--add-header", "Sec-Ch-Ua-Platform:\"Windows\"",
		"--add-header", "Sec-Fetch-Dest:document",
		"--add-header", "Sec-Fetch-Mode:navigate",
		"--add-header", "Sec-Fetch-Site:none",
		"--add-header", "Sec-Fetch-User:?1",
		"--add-header", "Upgrade-Insecure-Requests:1",
		"--extractor-retries", "15",
		"--fragment-retries", "20",
		"--retry-sleep", "exp=1:300",
		"--socket-timeout", "120",
		"--http-chunk-size", "1048576",
		"--sleep-interval", "2",
		"--max-sleep-interval", "10",
		"--geo-bypass",
		"--geo-bypass-country", "US",
		"--no-check-certificate",
		"--force-ipv4",
		"--verbose",
		"--prefer-free-formats",
		"--youtube-skip-dash-manifest",
		"--hls-prefer-native",
	}

	// Add FFmpeg location if available
	if ffmpegPath != "" {
		args = append(args, "--ffmpeg-location", ffmpegPath)
	}

	// Handle format selection - quality parameter is now a format_id
	if quality != "" && quality != "best" {
		// Check if quality is a format_id (numeric) or resolution (ends with 'p')
		if strings.HasSuffix(quality, "p") {
			// Legacy resolution format (e.g., "1080p")
			heightLimit := strings.TrimSuffix(quality, "p")
			formatString := fmt.Sprintf(
				"bestvideo[height=%s]+bestaudio[ext=m4a]/"+
					"bestvideo[height=%s]+bestaudio/"+
					"bestvideo[height<=%s]+bestaudio[ext=m4a]/"+
					"bestvideo[height<=%s]+bestaudio/"+
					"best[height<=%s]",
				heightLimit, heightLimit, heightLimit, heightLimit, heightLimit)
			args = append(args, "-f", formatString)
			log.Printf("Using resolution-based format string for %s: %s", quality, formatString)
		} else {
			// Modern format_id approach (e.g., "137", "270")
			formatString := fmt.Sprintf(
				"%s+bestaudio[ext=m4a]/"+
					"%s+bestaudio/"+
					"%s",
				quality, quality, quality)
			args = append(args, "-f", formatString)
			log.Printf("Using format_id-based selection for %s: %s", quality, formatString)
		}
	} else {
		// Force absolute best quality available
		formatString := "bestvideo[height>=1080]+bestaudio[ext=m4a]/" +
			"bestvideo[height>=1080]+bestaudio/" +
			"bestvideo[height>=720][fps>=30]+bestaudio[ext=m4a]/" +
			"bestvideo[height>=720][fps>=30]+bestaudio/" +
			"bestvideo+bestaudio[ext=m4a]/" +
			"bestvideo+bestaudio/" +
			"best[ext=mp4]/" +
			"best"
		args = append(args, "-f", formatString)
		log.Printf("Using best quality format string: %s", formatString)
	}

	args = append(args, url)

	// Get absolute path to yt-dlp
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("yt-dlp not found: %v", err))
		return
	}

	cmd := exec.Command(ytDlpPath, args...)
	log.Printf("=== DOWNLOADING WITH COMMAND ===")
	log.Printf("Path: %s", ytDlpPath)
	log.Printf("Full args: %v", args)
	log.Printf("Quality requested: %s", quality)
	log.Printf("URL: %s", url)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create stdout pipe: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create stderr pipe: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start yt-dlp: %v", err)
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to start yt-dlp: %v. Make sure yt-dlp.exe exists and is executable", err))
		return
	}

	// Parse progress output - multiple regex patterns for different yt-dlp outputs
	progressRegex1 := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%\s+of\s+.*?\s+at\s+(\S+)\s+ETA\s+(\S+)`)
	progressRegex2 := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%`)
	titleRegex := regexp.MustCompile(`\[download\] Destination: (.+)`)

	scanner := bufio.NewScanner(stdout)
	var title, filename string

	// Log stderr in background
	go func() {
		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			log.Printf("yt-dlp stderr: %s", line)
		}
	}()

	// Set initial status
	m.updateStatus(id, "downloading", 0, "", "", "Starting download...")

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("yt-dlp stdout: %s", line) // Debug logging

		// Extract title
		if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
			filename = filepath.Base(matches[1])
			title = strings.TrimSuffix(filename, filepath.Ext(filename))
			m.mu.Lock()
			if d, ok := m.downloads[id]; ok {
				d.Title = title
			}
			m.mu.Unlock()
		}

		// Extract progress - try multiple patterns
		if matches := progressRegex1.FindStringSubmatch(line); len(matches) > 3 {
			progress := parseFloat(matches[1])
			speed := matches[2]
			eta := matches[3]

			log.Printf("Progress update: %f%%, Speed: %s, ETA: %s", progress, speed, eta)
			m.updateStatus(id, "downloading", progress, speed, eta, "")
		} else if matches := progressRegex2.FindStringSubmatch(line); len(matches) > 1 {
			progress := parseFloat(matches[1])

			log.Printf("Simple progress update: %f%%", progress)
			m.updateStatus(id, "downloading", progress, "", "", "")
		}

		// Check for completion
		if strings.Contains(line, "[download] 100%") || strings.Contains(line, "has already been downloaded") {
			log.Printf("Download completed, processing...")
			m.updateStatus(id, "processing", 100, "", "", "Processing video...")
		}

		// Check for other status messages
		if strings.Contains(line, "[ffmpeg]") {
			m.updateStatus(id, "processing", 100, "", "", "Converting video...")
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("yt-dlp command failed: %v", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("yt-dlp exit code: %d", exitError.ExitCode())
			stderrOutput := string(exitError.Stderr)
			log.Printf("yt-dlp stderr: %s", stderrOutput)

			// Send more specific error message to user
			if strings.Contains(stderrOutput, "Video unavailable") {
				m.updateStatus(id, "error", 0, "", "", "Video is unavailable or private")
			} else if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "Forbidden") {
				m.updateStatus(id, "error", 0, "", "", "YouTube blocked the request. Try a different video or wait a moment.")
			} else if strings.Contains(stderrOutput, "Sign in") {
				m.updateStatus(id, "error", 0, "", "", "Video requires sign-in or is age-restricted")
			} else {
				m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Download failed: %s", stderrOutput))
			}
		} else {
			m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Download failed: %v", err))
		}
		return
	}

	// Success - download completed
	m.updateStatus(id, "processing", 100, "", "", "Opening File Explorer...")

	// Find the downloaded file in temp directory
	downloadedFile, err := m.findDownloadedFile(tempDir)
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Could not find downloaded file: %v", err))
		return
	}

	// Add resolution to filename
	newFileName := m.addResolutionToFilename(downloadedFile, quality)
	if newFileName != downloadedFile {
		if err := os.Rename(downloadedFile, newFileName); err != nil {
			log.Printf("Failed to rename file with resolution: %v", err)
		} else {
			downloadedFile = newFileName
		}
	}

	// Open File Picker for user to choose destination and filename
	finalPath, err := m.openFilePicker(downloadedFile)
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to save file: %v", err))
		return
	}

	m.updateStatus(id, "completed", 100, "", "", fmt.Sprintf("Saved as: %s", filepath.Base(finalPath)))
}

func (m *Manager) convertToMp3(id, url string) {
	m.mu.Lock()
	m.downloads[id] = &Download{
		ID:     id,
		URL:    url,
		Status: "starting",
	}
	m.mu.Unlock()

	// Clean the YouTube URL first (adapted from existing project)
	cleanURL, err := m.cleanYouTubeURL(url)
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Invalid YouTube URL: %v", err))
		return
	}

	// Create temp directory for MP3 conversion
	tempDir := filepath.Join(os.TempDir(), "ytmp3converter")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create temp directory: %v", err))
		return
	}

	outputPath := filepath.Join(tempDir, "%(title)s.%(ext)s")

	// Get FFmpeg path
	ffmpegPath, err := m.getFFmpegPath()
	if err != nil {
		log.Printf("Warning: FFmpeg not found, MP3 conversion may not work: %v", err)
	}

	// Prepare yt-dlp command for MP3 extraction
	args := []string{
		"--no-warnings",
		"--newline",
		"--progress",
		"-x", // Extract audio
		"--audio-format", "mp3",
		"--audio-quality", "0", // Best quality
		"-o", outputPath,
		"--embed-metadata",
		"--no-playlist",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"--referer", "https://www.youtube.com/",
		"--add-header", "Accept-Language:en-US,en;q=0.9",
		"--add-header", "Accept-Encoding:gzip, deflate, br, zstd",
		"--add-header", "Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"--add-header", "Cache-Control:no-cache",
		"--add-header", "Pragma:no-cache",
		"--add-header", "Sec-Ch-Ua:\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
		"--add-header", "Sec-Ch-Ua-Mobile:?0",
		"--add-header", "Sec-Ch-Ua-Platform:\"Windows\"",
		"--add-header", "Sec-Fetch-Dest:document",
		"--add-header", "Sec-Fetch-Mode:navigate",
		"--add-header", "Sec-Fetch-Site:none",
		"--add-header", "Sec-Fetch-User:?1",
		"--add-header", "Upgrade-Insecure-Requests:1",
		"--extractor-retries", "15",
		"--fragment-retries", "20",
		"--retry-sleep", "exp=1:300",
		"--socket-timeout", "120",
		"--http-chunk-size", "1048576",
		"--sleep-interval", "2",
		"--max-sleep-interval", "10",
		"--geo-bypass",
		"--geo-bypass-country", "US",
		"--no-check-certificate",
		"--force-ipv4",
		"--verbose",
		"--prefer-free-formats",
		"--youtube-skip-dash-manifest",
		"--hls-prefer-native",
	}

	// Add FFmpeg location if available
	if ffmpegPath != "" {
		args = append(args, "--ffmpeg-location", ffmpegPath)
	}

	args = append(args, cleanURL)

	// Get absolute path to yt-dlp
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("yt-dlp not found: %v", err))
		return
	}

	cmd := exec.Command(ytDlpPath, args...)
	log.Printf("=== CONVERTING TO MP3 ===")
	log.Printf("Path: %s", ytDlpPath)
	log.Printf("URL: %s -> %s", url, cleanURL)
	log.Printf("Temp dir: %s", tempDir)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create stdout pipe: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to create stderr pipe: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start yt-dlp for MP3: %v", err)
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to start MP3 conversion: %v", err))
		return
	}

	// Progress parsing (same as video download)
	progressRegex1 := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%\s+of\s+.*?\s+at\s+(\S+)\s+ETA\s+(\S+)`)
	progressRegex2 := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%`)
	titleRegex := regexp.MustCompile(`\[download\] Destination: (.+)`)

	scanner := bufio.NewScanner(stdout)
	var title, filename string

	// Log stderr in background
	go func() {
		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			log.Printf("yt-dlp MP3 stderr: %s", line)
		}
	}()

	// Set initial status
	m.updateStatus(id, "converting", 0, "", "", "Starting MP3 conversion...")

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("yt-dlp MP3 stdout: %s", line)

		// Extract title
		if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
			filename = filepath.Base(matches[1])
			title = strings.TrimSuffix(filename, filepath.Ext(filename))
			m.mu.Lock()
			if d, ok := m.downloads[id]; ok {
				d.Title = title
			}
			m.mu.Unlock()
		}

		// Extract progress
		if matches := progressRegex1.FindStringSubmatch(line); len(matches) > 3 {
			progress := parseFloat(matches[1])
			speed := matches[2]
			eta := matches[3]

			log.Printf("MP3 Progress: %f%%, Speed: %s, ETA: %s", progress, speed, eta)
			m.updateStatus(id, "converting", progress, speed, eta, "")
		} else if matches := progressRegex2.FindStringSubmatch(line); len(matches) > 1 {
			progress := parseFloat(matches[1])
			log.Printf("MP3 Simple progress: %f%%", progress)
			m.updateStatus(id, "converting", progress, "", "", "")
		}

		// Check for completion
		if strings.Contains(line, "[download] 100%") || strings.Contains(line, "has already been downloaded") {
			log.Printf("MP3 conversion completed, processing...")
			m.updateStatus(id, "processing", 100, "", "", "Processing MP3...")
		}

		// Check for FFmpeg processing
		if strings.Contains(line, "[ffmpeg]") {
			m.updateStatus(id, "processing", 100, "", "", "Converting to MP3...")
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("yt-dlp MP3 command failed: %v", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("yt-dlp MP3 exit code: %d", exitError.ExitCode())
			stderrOutput := string(exitError.Stderr)
			log.Printf("yt-dlp MP3 stderr: %s", stderrOutput)

			// Send specific error messages
			if strings.Contains(stderrOutput, "Video unavailable") {
				m.updateStatus(id, "error", 0, "", "", "Video is unavailable or private")
			} else if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "Forbidden") {
				m.updateStatus(id, "error", 0, "", "", "YouTube blocked the request. Try a different video or wait a moment.")
			} else if strings.Contains(stderrOutput, "Sign in") {
				m.updateStatus(id, "error", 0, "", "", "Video requires sign-in or is age-restricted")
			} else {
				m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("MP3 conversion failed: %s", stderrOutput))
			}
		} else {
			m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("MP3 conversion failed: %v", err))
		}
		return
	}

	// Success - find the converted MP3 file
	m.updateStatus(id, "processing", 100, "", "", "Opening File Picker...")

	convertedFile, err := m.findDownloadedFile(tempDir)
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Could not find converted MP3 file: %v", err))
		return
	}

	// Open File Picker for user to choose destination
	finalPath, err := m.openFilePicker(convertedFile)
	if err != nil {
		m.updateStatus(id, "error", 0, "", "", fmt.Sprintf("Failed to save MP3 file: %v", err))
		return
	}

	m.updateStatus(id, "completed", 100, "", "", fmt.Sprintf("MP3 saved as: %s", filepath.Base(finalPath)))
}

func (m *Manager) cleanYouTubeURL(rawURL string) (string, error) {
	// Parse and clean the YouTube URL (adapted from existing project)
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Handle youtu.be format
	if strings.Contains(u.Host, "youtu.be") {
		videoID := strings.TrimPrefix(u.Path, "/")
		return "https://www.youtube.com/watch?v=" + videoID, nil
	}

	// Handle youtube.com format
	q := u.Query()
	videoID := q.Get("v")
	if videoID == "" {
		return "", fmt.Errorf("missing v parameter in URL")
	}

	return "https://www.youtube.com/watch?v=" + videoID, nil
}

func (m *Manager) updateStatus(id, status string, progress float64, speed, eta, message string) {
	m.mu.Lock()
	if download, ok := m.downloads[id]; ok {
		download.Status = status
		download.Progress = progress
		download.Speed = speed
		download.ETA = eta
	}
	m.mu.Unlock()

	update := models.ProgressUpdate{
		ID:       id,
		Progress: progress,
		Speed:    speed,
		ETA:      eta,
		Status:   status,
		Message:  message,
	}

	log.Printf("Broadcasting update: %+v", update)
	m.broadcast(update)
}

func (m *Manager) broadcast(update models.ProgressUpdate) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	log.Printf("Broadcasting to %d subscribers", len(m.subscribers))
	for i, ch := range m.subscribers {
		select {
		case ch <- update:
			log.Printf("Sent update to subscriber %d", i)
		default:
			log.Printf("Subscriber %d channel full, skipping", i)
		}
	}
}

func (m *Manager) SubscribeToUpdates() <-chan models.ProgressUpdate {
	ch := make(chan models.ProgressUpdate, 100)
	m.mu.Lock()
	m.subscribers = append(m.subscribers, ch)
	m.mu.Unlock()
	return ch
}

func (m *Manager) ParseYouTubeURL(inputURL string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Check if it's a YouTube URL
	host := strings.ToLower(parsedURL.Host)
	if !strings.Contains(host, "youtube.com") && !strings.Contains(host, "youtu.be") {
		return "", fmt.Errorf("not a YouTube URL")
	}

	var videoID string

	// Handle different YouTube URL formats
	if strings.Contains(host, "youtu.be") {
		// Format: https://youtu.be/VIDEO_ID
		videoID = strings.TrimPrefix(parsedURL.Path, "/")
	} else {
		// Format: https://www.youtube.com/watch?v=VIDEO_ID
		videoID = parsedURL.Query().Get("v")
		if videoID == "" {
			// Try to extract from embed format
			if strings.Contains(parsedURL.Path, "/embed/") {
				parts := strings.Split(parsedURL.Path, "/embed/")
				if len(parts) > 1 {
					videoID = parts[1]
				}
			}
		}
	}

	if videoID == "" {
		return "", fmt.Errorf("could not extract video ID from URL")
	}

	// Return clean URL without playlist or other parameters
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), nil
}

func (m *Manager) GetVideoInfo(url string) (*models.VideoInfo, error) {
	// First parse the URL
	parsedURL, err := m.ParseYouTubeURL(url)
	if err != nil {
		return nil, err
	}

	// Get absolute path to yt-dlp
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp not found: %v", err)
	}

	// Use yt-dlp to get video info with maximum anti-blocking measures
	cmd := exec.Command(ytDlpPath,
		"-j",
		"--no-warnings",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"--referer", "https://www.youtube.com/",
		"--add-header", "Accept-Language:en-US,en;q=0.9",
		"--add-header", "Accept-Encoding:gzip, deflate, br, zstd",
		"--add-header", "Sec-Ch-Ua:\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
		"--add-header", "Sec-Ch-Ua-Mobile:?0",
		"--add-header", "Sec-Ch-Ua-Platform:\"Windows\"",
		"--geo-bypass",
		"--geo-bypass-country", "US",
		"--extractor-retries", "10",
		"--fragment-retries", "10",
		"--retry-sleep", "exp=1:120",
		"--no-check-certificate",
		"--force-ipv4",
		"--prefer-free-formats",
		"--youtube-skip-dash-manifest",
		parsedURL)
	log.Printf("Getting video info with command: %s -j --no-warnings %s", ytDlpPath, parsedURL)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("yt-dlp video info command failed: %v", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			stderrOutput := string(exitError.Stderr)
			log.Printf("yt-dlp stderr: %s", stderrOutput)

			// Handle specific error cases
			if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "Forbidden") {
				return nil, fmt.Errorf("video is restricted or geo-blocked. YouTube blocked access to this video")
			} else if strings.Contains(stderrOutput, "fragment") && strings.Contains(stderrOutput, "not found") {
				return nil, fmt.Errorf("video fragments are unavailable. This video may be corrupted or restricted")
			} else if strings.Contains(stderrOutput, "Private video") {
				return nil, fmt.Errorf("this is a private video")
			} else if strings.Contains(stderrOutput, "Video unavailable") {
				return nil, fmt.Errorf("video is unavailable or has been removed")
			}

			return nil, fmt.Errorf("failed to get video info: %s", stderrOutput)
		}
		return nil, fmt.Errorf("failed to get video info: %v", err)
	}

	// Parse JSON output
	log.Printf("=== RAW VIDEO INFO OUTPUT ===")
	log.Printf("Output length: %d bytes", len(output))
	if len(output) > 1000 {
		log.Printf("First 1000 chars: %s", string(output[:1000]))
	} else {
		log.Printf("Full output: %s", string(output))
	}

	var rawInfo map[string]interface{}
	if err := json.Unmarshal(output, &rawInfo); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %v", err)
	}

	videoInfo := &models.VideoInfo{
		ParsedURL: parsedURL,
	}

	// Extract basic info
	if title, ok := rawInfo["title"].(string); ok {
		videoInfo.Title = title
	}

	if duration, ok := rawInfo["duration"].(float64); ok {
		minutes := int(duration) / 60
		seconds := int(duration) % 60
		videoInfo.Duration = fmt.Sprintf("%d:%02d", minutes, seconds)
	}

	if thumbnail, ok := rawInfo["thumbnail"].(string); ok {
		videoInfo.Thumbnail = thumbnail
	}

	// Extract formats
	if formats, ok := rawInfo["formats"].([]interface{}); ok {
		log.Printf("=== AVAILABLE FORMATS ===")
		log.Printf("Total formats found: %d", len(formats))

		qualityMap := make(map[string]models.VideoFormat)

		for i, f := range formats {
			if format, ok := f.(map[string]interface{}); ok {
				// Log format details for debugging
				if i < 10 { // Only log first 10 for brevity
					log.Printf("Format %d: height=%v, vcodec=%v, acodec=%v, ext=%v, format_id=%v, tbr=%v",
						i, format["height"], format["vcodec"], format["acodec"], format["ext"], format["format_id"], format["tbr"])
				}

				// Skip audio-only formats
				if vcodec, ok := format["vcodec"].(string); ok && vcodec == "none" {
					continue
				}

				// Get resolution
				resolution := "Unknown"
				if height, ok := format["height"].(float64); ok {
					resolution = fmt.Sprintf("%dp", int(height))
					log.Printf("Found video format: %s (format_id: %v)", resolution, format["format_id"])
				}

				// Skip if we already have this resolution
				if _, exists := qualityMap[resolution]; exists {
					continue
				}

				videoFormat := models.VideoFormat{
					Resolution: resolution,
				}

				if formatID, ok := format["format_id"].(string); ok {
					videoFormat.FormatID = formatID
				}

				if ext, ok := format["ext"].(string); ok {
					videoFormat.Extension = ext
				}

				if filesize, ok := format["filesize"].(float64); ok {
					if filesize > 0 {
						videoFormat.FileSize = formatFileSize(int64(filesize))
					} else if filesizeApprox, ok := format["filesize_approx"].(float64); ok {
						videoFormat.FileSize = "~" + formatFileSize(int64(filesizeApprox))
					}
				}

				// Store best format for each resolution
				qualityMap[resolution] = videoFormat
			}
		}

		// Convert map to slice
		for _, format := range qualityMap {
			videoInfo.Formats = append(videoInfo.Formats, format)
		}

		// Sort formats by resolution (highest first)
		sortFormatsByResolution(videoInfo.Formats)
	}

	return videoInfo, nil
}

func (m *Manager) openFilePicker(sourceFile string) (string, error) {
	filename := filepath.Base(sourceFile)

	switch runtime.GOOS {
	case "windows":
		// Use PowerShell to show SaveFileDialog
		psScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
$saveFileDialog = New-Object System.Windows.Forms.SaveFileDialog
$saveFileDialog.Filter = "Video Files|*.mp4;*.mkv;*.avi;*.mov;*.wmv;*.flv;*.webm|All Files|*.*"
$saveFileDialog.FileName = "%s"
$saveFileDialog.Title = "Save Video As"
$saveFileDialog.DefaultExt = "%s"

if ($saveFileDialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $saveFileDialog.FileName
} else {
    Write-Output "CANCELLED"
}`, filename, filepath.Ext(filename))

		cmd := exec.Command("powershell", "-Command", psScript)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("PowerShell SaveFileDialog failed: %v", err)
			// Fallback to Downloads folder
			return m.fallbackToDownloads(sourceFile)
		}

		selectedPath := strings.TrimSpace(string(output))
		if selectedPath == "CANCELLED" || selectedPath == "" {
			// User cancelled - clean up temp file and return error
			os.Remove(sourceFile)
			return "", fmt.Errorf("save cancelled by user")
		}

		// Copy file to selected location
		if err := m.copyFile(sourceFile, selectedPath); err != nil {
			return "", fmt.Errorf("failed to save file: %v", err)
		}

		// Clean up temp file
		os.Remove(sourceFile)
		return selectedPath, nil

	case "darwin":
		// macOS doesn't have a simple command-line save dialog, fallback to Downloads
		log.Printf("File picker not implemented for macOS, using Downloads folder")
		return m.fallbackToDownloads(sourceFile)

	case "linux":
		// Linux doesn't have a universal save dialog, fallback to Downloads
		log.Printf("File picker not implemented for Linux, using Downloads folder")
		return m.fallbackToDownloads(sourceFile)

	default:
		log.Printf("Unsupported OS: %s, using Downloads folder", runtime.GOOS)
		return m.fallbackToDownloads(sourceFile)
	}
}

func (m *Manager) fallbackToDownloads(sourceFile string) (string, error) {
	// Create downloads folder as default
	downloadsDir, err := os.UserHomeDir()
	if err != nil {
		downloadsDir = "."
	}
	downloadsDir = filepath.Join(downloadsDir, "Downloads")
	os.MkdirAll(downloadsDir, 0755)

	// Get unique filename to handle duplicates
	filename := filepath.Base(sourceFile)
	uniquePath := m.getUniqueFilePath(downloadsDir, filename)

	if err := m.copyFile(sourceFile, uniquePath); err != nil {
		return "", fmt.Errorf("failed to copy to downloads folder: %v", err)
	}

	// Clean up temp file
	os.Remove(sourceFile)

	return uniquePath, nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func sortFormatsByResolution(formats []models.VideoFormat) {
	// Simple sort by resolution (highest first)
	for i := 0; i < len(formats); i++ {
		for j := i + 1; j < len(formats); j++ {
			res1 := extractResolutionNumber(formats[i].Resolution)
			res2 := extractResolutionNumber(formats[j].Resolution)
			if res1 < res2 {
				formats[i], formats[j] = formats[j], formats[i]
			}
		}
	}
}

func extractResolutionNumber(resolution string) int {
	// Extract number from resolution like "1080p" -> 1080
	var num int
	fmt.Sscanf(resolution, "%dp", &num)
	return num
}
