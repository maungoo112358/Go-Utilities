package downloader

import (
	"Go-Utilities/internal/consts"
	"Go-Utilities/internal/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ExecuteDownload(url, quality string, progressCallback ProgressCallback) (*YtDlpResult, error) {
	tempDir := filepath.Join(os.TempDir(), consts.TEMP_DIR)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf(consts.ERR_CREATE_TEMP_DIR, err)
	}

	outputPath := filepath.Join(tempDir, consts.YT_DLP_OUTPUT_FORMAT)

	ffmpegPath, err := getFFmpegPath()
	if err != nil {
		log.Printf(consts.WARNING_FFMPEG_NOT_FOUND, err)
	}

	args := []string{"-o", outputPath}
	args = append(args, consts.YT_DLP_DOWNLOAD_ARGS...)

	if ffmpegPath != "" {
		args = append(args, consts.FFMPEG_LOCATION_FLAG, ffmpegPath)
	}

	if quality != "" && quality != "best" {
		if strings.HasSuffix(quality, "p") {
			heightLimit := strings.TrimSuffix(quality, "p")
			formatString := fmt.Sprintf(
				"bestvideo[height<=%s][ext=mp4]+bestaudio[ext=m4a]/"+
					"bestvideo[height<=%s]+bestaudio/"+
					"best[height<=%s]",
				heightLimit, heightLimit, heightLimit)
			args = append(args, "-f", formatString)
		} else {
			formatString := fmt.Sprintf(
				"(%s+bestaudio[ext=m4a])/"+
					"(%s+bestaudio)/"+
					"%s/best",
				quality, quality, quality)
			args = append(args, "-f", formatString)
		}
	} else {
		formatString := "bestvideo[height>=1080]+bestaudio[ext=m4a]/" +
			"bestvideo[height>=1080]+bestaudio/" +
			"bestvideo[height>=720][fps>=30]+bestaudio[ext=m4a]/" +
			"bestvideo[height>=720][fps>=30]+bestaudio/" +
			"bestvideo[height>=720]+bestaudio[ext=m4a]/" +
			"bestvideo[height>=720]+bestaudio/" +
			"best[height>=720]/best"
		args = append(args, "-f", formatString)
	}

	args = append(args, url)

	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_START_YT_DLP_DOWNLOAD, err)
	}

	log.Printf(consts.LOG_DOWNLOADING_COMMAND)
	log.Printf(consts.LOG_PATH, ytDlpPath)
	log.Printf(consts.LOG_FULL_ARGS, args)
	log.Printf(consts.LOG_QUALITY_REQUESTED, quality)
	log.Printf(consts.LOG_URL, url)

	cmd := exec.Command(ytDlpPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_CREATE_STDOUT_PIPE, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_CREATE_STDERR_PIPE, err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(consts.ERR_START_YT_DLP_EXE, err)
	}

	progressRegex1 := regexp.MustCompile(consts.YT_DLP_PROGRESS_REGEX_WITH_SPEED)
	progressRegex2 := regexp.MustCompile(consts.YT_DLP_PROGRESS_REGEX_SIMPLE)
	titleRegex := regexp.MustCompile(consts.YT_DLP_TITLE_REGEX)

	scanner := bufio.NewScanner(stdout)
	var title, filename string

	go func() {
		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			log.Printf(consts.LOG_YT_DLP_STDERR, line)
		}
	}()

	for scanner.Scan() {
		line := scanner.Text()

		if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
			filename = filepath.Base(matches[1])
			title = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		if matches := progressRegex1.FindStringSubmatch(line); len(matches) > 3 {
			progress := parseFloat(matches[1])
			speed := matches[2]
			eta := matches[3]
			if progressCallback != nil {
				progressCallback(progress, speed, eta, consts.MSG_STARTING_DOWNLOAD)
			}
		} else if matches := progressRegex2.FindStringSubmatch(line); len(matches) > 1 {
			progress := parseFloat(matches[1])
			if progressCallback != nil {
				progressCallback(progress, "", "", consts.MSG_STARTING_DOWNLOAD)
			}
		}

		if strings.Contains(line, "[download] 100%") || strings.Contains(line, "has already been downloaded") {
			if progressCallback != nil {
				progressCallback(100, "", "", consts.MSG_DOWNLOAD_COMPLETE)
			}
		}

		if strings.Contains(line, "[ffmpeg]") {
			if progressCallback != nil {
				progressCallback(100, "", "", consts.MSG_CONVERTING_VIDEO)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			stderrOutput := string(exitError.Stderr)
			if strings.Contains(stderrOutput, "Video unavailable") {
				return nil, fmt.Errorf(consts.ERR_VIDEO_UNAVAILABLE)
			} else if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "Forbidden") {
				return nil, fmt.Errorf(consts.ERR_YOUTUBE_BLOCKED)
			} else if strings.Contains(stderrOutput, "Sign in") {
				return nil, fmt.Errorf(consts.ERR_VIDEO_RESTRICTED)
			}
			return nil, fmt.Errorf(consts.ERR_DOWNLOAD_FAILED, stderrOutput)
		}
		return nil, fmt.Errorf(consts.ERR_DOWNLOAD_FAILED, err.Error())
	}

	downloadedFile, err := findDownloadedFile(tempDir)
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_FIND_DOWNLOADED_FILE, err)
	}

	return &YtDlpResult{
		Title:    title,
		FilePath: downloadedFile,
		TempDir:  tempDir,
		Success:  true,
	}, nil
}

func GetVideoInfo(url string) (*models.VideoInfo, error) {
	parsedURL, err := ParseYouTubeURL(url)
	if err != nil {
		return nil, err
	}

	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_START_YT_DLP_INFO, err)
	}

	args := append(consts.YT_DLP_INFO_ARGS, parsedURL)

	log.Printf(consts.LOG_GETTING_VIDEO_INFO, ytDlpPath, parsedURL)

	cmd := exec.Command(ytDlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Printf(consts.ERR_YT_DLP_INFO_FAILED, err)
		if exitError, ok := err.(*exec.ExitError); ok {
			stderrOutput := string(exitError.Stderr)
			log.Printf(consts.LOG_YT_DLP_STDERR, stderrOutput)

			if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "Forbidden") {
				return nil, fmt.Errorf(consts.ERR_VIDEO_RESTRICTED_GEO)
			} else if strings.Contains(stderrOutput, "fragment") && strings.Contains(stderrOutput, "not found") {
				return nil, fmt.Errorf(consts.ERR_VIDEO_FRAGMENTS)
			} else if strings.Contains(stderrOutput, "Private video") {
				return nil, fmt.Errorf(consts.ERR_VIDEO_PRIVATE)
			} else if strings.Contains(stderrOutput, "Video unavailable") {
				return nil, fmt.Errorf(consts.ERR_VIDEO_REMOVED)
			}

			return nil, fmt.Errorf(consts.ERR_GET_VIDEO_INFO, stderrOutput)
		}
		return nil, fmt.Errorf(consts.ERR_GET_VIDEO_INFO_2, err)
	}

	log.Printf(consts.LOG_RAW_VIDEO_INFO)
	log.Printf(consts.LOG_OUTPUT_LENGTH, len(output))
	if len(output) > 1000 {
		log.Printf(consts.LOG_FIRST_1000_CHARS, string(output[:1000]))
	} else {
		log.Printf(consts.LOG_FULL_OUTPUT, string(output))
	}

	var rawInfo map[string]interface{}
	if err := json.Unmarshal(output, &rawInfo); err != nil {
		return nil, fmt.Errorf(consts.ERR_PARSE_VIDEO_INFO, err)
	}

	videoInfo := &models.VideoInfo{
		Formats:   []models.VideoFormat{},
		ParsedURL: parsedURL,
	}

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

	if formats, ok := rawInfo["formats"].([]interface{}); ok {
		log.Printf(consts.LOG_AVAILABLE_FORMATS)
		log.Printf(consts.LOG_TOTAL_FORMATS, len(formats))

		qualityMap := make(map[string]models.VideoFormat)

		for i, f := range formats {
			if format, ok := f.(map[string]interface{}); ok {
				if i < 10 {
					log.Printf(consts.LOG_FORMAT_DETAILS,
						i, format["height"], format["vcodec"], format["acodec"], format["ext"], format["format_id"], format["tbr"])
				}

				if vcodec, ok := format["vcodec"].(string); ok && vcodec == "none" {
					continue
				}

				var videoFormat models.VideoFormat
				resolution := "unknown"

				if height, ok := format["height"].(float64); ok {
					resolution = fmt.Sprintf("%dp", int(height))
					log.Printf(consts.LOG_FOUND_VIDEO_FORMAT, resolution, format["format_id"])
				}

				if _, exists := qualityMap[resolution]; exists {
					continue
				}

				videoFormat.Resolution = resolution

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

				qualityMap[resolution] = videoFormat
			}
		}

		for _, format := range qualityMap {
			videoInfo.Formats = append(videoInfo.Formats, format)
		}

		sortFormatsByResolution(videoInfo.Formats)
	}

	return videoInfo, nil
}

func sortFormatsByResolution(formats []models.VideoFormat) {
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
	var num int
	fmt.Sscanf(resolution, "%dp", &num)
	return num
}