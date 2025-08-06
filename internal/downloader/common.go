package downloader

import (
	"Go-Utilities/internal/consts"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type YtDlpResult struct {
	Title    string
	FilePath string
	TempDir  string
	Success  bool
	Error    string
}

type ProgressCallback func(progress float64, speed, eta, message string)

func getYtDlpPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	ytDlpPath := filepath.Join(wd, consts.DEPENDENCIES_DIR, consts.YT_DLP_EXE_NAME)
	if _, err := os.Stat(ytDlpPath); os.IsNotExist(err) {
		return "", fmt.Errorf(consts.ERR_YT_DLP_NOT_FOUND, ytDlpPath)
	}
	return ytDlpPath, nil
}

func getFFmpegPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	ffmpegPath := filepath.Join(wd, consts.DEPENDENCIES_DIR, consts.FFMPEG_EXE_NAME)
	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
		return "", fmt.Errorf(consts.ERR_FFMPEG_NOT_FOUND, ffmpegPath)
	}
	return ffmpegPath, nil
}

func TestYtDlp() error {
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(ytDlpPath, consts.YT_DLP_VERSION_FLAG)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf(consts.ERR_YT_DLP_TEST_FAILED, err)
	}
	version := strings.TrimSpace(string(output))
	log.Printf(consts.LOG_YT_DLP_VERSION, version)
	if len(version) > 0 && version < consts.MIN_YTDLP_YEAR {
		log.Printf(consts.WARNING_YT_DLP_OUTDATED)
	}
	return nil
}

func findDownloadedFile(dir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			name := strings.ToLower(filepath.Base(file))
			if !strings.Contains(name, consts.FRAGMENT_EXT) && !strings.Contains(name, consts.TEMP_EXT) {
				return file, nil
			}
		}
	}

	return "", fmt.Errorf(consts.ERR_NO_VIDEO_FILE, dir)
}

func extractVideoID(parsedURL *url.URL) string {
	host := strings.ToLower(parsedURL.Host)
	
	if strings.Contains(host, consts.YOUTU_BE_DOMAIN) {
		return strings.TrimPrefix(parsedURL.Path, consts.URL_PATH_SEPARATOR)
	}
	
	videoID := parsedURL.Query().Get(consts.URL_VIDEO_PARAM)
	if videoID == "" {
		if strings.Contains(parsedURL.Path, consts.YOUTUBE_EMBED_PATH) {
			parts := strings.Split(parsedURL.Path, consts.YOUTUBE_EMBED_PATH)
			if len(parts) > 1 {
				videoID = strings.Split(parts[1], consts.URL_PATH_SEPARATOR)[0]
			}
		}
	}
	
	return videoID
}

func ParseYouTubeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf(consts.ERR_INVALID_URL, err)
	}

	host := strings.ToLower(parsedURL.Host)
	if !strings.Contains(host, consts.YOUTUBE_DOMAIN) && !strings.Contains(host, consts.YOUTU_BE_DOMAIN) {
		return "", fmt.Errorf(consts.ERR_NOT_YOUTUBE_URL)
	}

	videoID := extractVideoID(parsedURL)
	if videoID == "" {
		return "", fmt.Errorf(consts.ERR_EXTRACT_VIDEO_ID)
	}

	return fmt.Sprintf(consts.YOUTUBE_WATCH_URL, videoID), nil
}

func cleanYouTubeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	videoID := extractVideoID(parsedURL)
	if videoID == "" {
		return "", fmt.Errorf(consts.ERR_MISSING_V_PARAM)
	}

	return fmt.Sprintf(consts.YOUTUBE_WATCH_URL, videoID), nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func formatFileSize(bytes int64) string {
	const unit = consts.BYTES_UNIT
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), consts.FILE_SIZE_UNITS[exp])
}