package downloader

import (
	"Go-Utilities/internal/consts"
	"Go-Utilities/internal/models"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ExecuteDownload(url, quality string, progressCallback ProgressCallback) (*YtDlpResult, error) {
	tempDir, err := prepareDownloadEnvironment()
	if err != nil {
		return nil, err
	}

	args, err := buildDownloadCommand(tempDir, url, quality)
	if err != nil {
		return nil, err
	}

	title, err := executeDownloadProcess(args, url, quality, progressCallback)
	if err != nil {
		return nil, err
	}

	return locateDownloadResult(tempDir, title)
}

func prepareDownloadEnvironment() (string, error) {
	tempDir := filepath.Join(os.TempDir(), consts.TEMP_DIR)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf(consts.ERR_CREATE_TEMP_DIR, err)
	}
	return tempDir, nil
}

func buildDownloadCommand(tempDir, url, quality string) ([]string, error) {
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

	args = appendQualityFormat(args, quality)
	args = append(args, url)

	return args, nil
}

func appendQualityFormat(args []string, quality string) []string {
	if quality != "" && quality != consts.BEST_QUALITY {
		if strings.HasSuffix(quality, consts.QUALITY_SUFFIX) {
			heightLimit := strings.TrimSuffix(quality, consts.QUALITY_SUFFIX)
			formatString := fmt.Sprintf(consts.QUALITY_HEIGHT_FORMAT, heightLimit, heightLimit, heightLimit)
			args = append(args, consts.FORMAT_FLAG, formatString)
		} else {
			formatString := fmt.Sprintf(consts.QUALITY_CUSTOM_FORMAT, quality, quality, quality)
			args = append(args, consts.FORMAT_FLAG, formatString)
		}
	} else {
		args = append(args, consts.FORMAT_FLAG, consts.QUALITY_BEST_FORMAT)
	}
	return args
}

func executeDownloadProcess(args []string, url, quality string, progressCallback ProgressCallback) (string, error) {
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return "", fmt.Errorf(consts.ERR_START_YT_DLP_DOWNLOAD, err)
	}

	log.Printf(consts.LOG_DOWNLOAD_COMMAND_INFO, ytDlpPath, args, quality, url)

	cmd := exec.Command(ytDlpPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf(consts.ERR_CREATE_STDOUT_PIPE, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf(consts.ERR_CREATE_STDERR_PIPE, err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf(consts.ERR_START_YT_DLP_EXE, err)
	}

	go handleDownloadStderrOutput(stderr)

	title := handleDownloadProcessOutput(stdout, progressCallback)

	if err := cmd.Wait(); err != nil {
		return "", validateDownloadResult(err)
	}

	return title, nil
}

func handleDownloadStderrOutput(stderr io.ReadCloser) {
	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		line := stderrScanner.Text()
		log.Printf(consts.LOG_YT_DLP_STDERR, line)
	}
}

func handleDownloadProcessOutput(stdout io.ReadCloser, progressCallback ProgressCallback) string {
	progressRegex1 := regexp.MustCompile(consts.YT_DLP_PROGRESS_REGEX_WITH_SPEED)
	progressRegex2 := regexp.MustCompile(consts.YT_DLP_PROGRESS_REGEX_SIMPLE)
	titleRegex := regexp.MustCompile(consts.YT_DLP_TITLE_REGEX)

	scanner := bufio.NewScanner(stdout)
	var title, filename string

	for scanner.Scan() {
		line := scanner.Text()

		if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
			filename = filepath.Base(matches[1])
			title = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		parseDownloadProgress(line, progressRegex1, progressRegex2, progressCallback)
	}

	return title
}

func parseDownloadProgress(line string, progressRegex1, progressRegex2 *regexp.Regexp, progressCallback ProgressCallback) {
	if progressCallback == nil {
		return
	}

	if matches := progressRegex1.FindStringSubmatch(line); len(matches) > 3 {
		progress := parseFloat(matches[1])
		speed := matches[2]
		eta := matches[3]
		progressCallback(progress, speed, eta, consts.MSG_STARTING_DOWNLOAD)
	} else if matches := progressRegex2.FindStringSubmatch(line); len(matches) > 1 {
		progress := parseFloat(matches[1])
		progressCallback(progress, "", "", consts.MSG_STARTING_DOWNLOAD)
	}

	if strings.Contains(line, consts.YT_DLP_DOWNLOAD_100_PERCENT) || strings.Contains(line, consts.YT_DLP_ALREADY_DOWNLOADED) {
		progressCallback(100, "", "", consts.MSG_DOWNLOAD_COMPLETE)
	}

	if strings.Contains(line, consts.YT_DLP_FFMPEG_TAG) {
		progressCallback(100, "", "", consts.MSG_CONVERTING_VIDEO)
	}
}

func validateDownloadResult(err error) error {
	if exitError, ok := err.(*exec.ExitError); ok {
		stderrOutput := string(exitError.Stderr)
		if strings.Contains(stderrOutput, consts.YT_DLP_VIDEO_UNAVAILABLE) {
			return fmt.Errorf(consts.ERR_VIDEO_UNAVAILABLE)
		} else if strings.Contains(stderrOutput, consts.YT_DLP_FORBIDDEN_403) || strings.Contains(stderrOutput, consts.YT_DLP_FORBIDDEN_TEXT) {
			return fmt.Errorf(consts.ERR_YOUTUBE_BLOCKED)
		} else if strings.Contains(stderrOutput, consts.YT_DLP_SIGN_IN_REQUIRED) {
			return fmt.Errorf(consts.ERR_VIDEO_RESTRICTED)
		}
		return fmt.Errorf(consts.ERR_DOWNLOAD_FAILED, stderrOutput)
	}
	return fmt.Errorf(consts.ERR_DOWNLOAD_FAILED, err.Error())
}

func locateDownloadResult(tempDir, title string) (*YtDlpResult, error) {
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

	rawOutput, err := executeVideoInfoCommand(parsedURL)
	if err != nil {
		return nil, err
	}

	rawInfo, err := parseVideoInfoJSON(rawOutput)
	if err != nil {
		return nil, err
	}

	videoInfo := extractBasicVideoInfo(rawInfo, parsedURL)
	processVideoFormats(rawInfo, videoInfo)

	return videoInfo, nil
}

func executeVideoInfoCommand(parsedURL string) ([]byte, error) {
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_START_YT_DLP_INFO, err)
	}

	args := append(consts.YT_DLP_INFO_ARGS, parsedURL)
	log.Printf(consts.LOG_GETTING_VIDEO_INFO, ytDlpPath, parsedURL)

	cmd := exec.Command(ytDlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, validateVideoInfoError(err)
	}

	return output, nil
}

func validateVideoInfoError(err error) error {
	if exitError, ok := err.(*exec.ExitError); ok {
		stderrOutput := string(exitError.Stderr)
		log.Printf(consts.LOG_YT_DLP_INFO_FAILED_STDERR, err, stderrOutput)

		if strings.Contains(stderrOutput, consts.YT_DLP_FORBIDDEN_403) || strings.Contains(stderrOutput, consts.YT_DLP_FORBIDDEN_TEXT) {
			return fmt.Errorf(consts.ERR_VIDEO_RESTRICTED_GEO)
		} else if strings.Contains(stderrOutput, consts.YT_DLP_FRAGMENT_TEXT) && strings.Contains(stderrOutput, consts.YT_DLP_NOT_FOUND_TEXT) {
			return fmt.Errorf(consts.ERR_VIDEO_FRAGMENTS)
		} else if strings.Contains(stderrOutput, consts.YT_DLP_PRIVATE_VIDEO_TEXT) {
			return fmt.Errorf(consts.ERR_VIDEO_PRIVATE)
		} else if strings.Contains(stderrOutput, consts.YT_DLP_VIDEO_UNAVAILABLE) {
			return fmt.Errorf(consts.ERR_VIDEO_REMOVED)
		}

		return fmt.Errorf(consts.ERR_GET_VIDEO_INFO, stderrOutput)
	}
	return fmt.Errorf(consts.ERR_GET_VIDEO_INFO_2, err)
}

func parseVideoInfoJSON(output []byte) (map[string]interface{}, error) {
	log.Printf(consts.LOG_RAW_VIDEO_INFO_LENGTH, len(output))
	if len(output) > 1000 {
		log.Printf(consts.LOG_FIRST_1000_CHARS, string(output[:1000]))
	} else {
		log.Printf(consts.LOG_FULL_OUTPUT, string(output))
	}

	var rawInfo map[string]interface{}
	if err := json.Unmarshal(output, &rawInfo); err != nil {
		return nil, fmt.Errorf(consts.ERR_PARSE_VIDEO_INFO, err)
	}

	return rawInfo, nil
}

func extractBasicVideoInfo(rawInfo map[string]interface{}, parsedURL string) *models.VideoInfo {
	videoInfo := &models.VideoInfo{
		Formats:   []models.VideoFormat{},
		ParsedURL: parsedURL,
	}

	if title, ok := rawInfo[consts.JSON_TITLE].(string); ok {
		videoInfo.Title = title
	}

	if duration, ok := rawInfo[consts.JSON_DURATION].(float64); ok {
		minutes := int(duration) / 60
		seconds := int(duration) % 60
		videoInfo.Duration = fmt.Sprintf(consts.DURATION_FORMAT, minutes, seconds)
	}

	if thumbnail, ok := rawInfo[consts.JSON_THUMBNAIL].(string); ok {
		videoInfo.Thumbnail = thumbnail
	}

	return videoInfo
}

func processVideoFormats(rawInfo map[string]interface{}, videoInfo *models.VideoInfo) {
	formats, ok := rawInfo[consts.JSON_FORMATS].([]interface{})
	if !ok {
		return
	}

	log.Printf(consts.LOG_AVAILABLE_FORMATS_TOTAL, len(formats))

	qualityMap := buildQualityMap(formats)
	
	for _, format := range qualityMap {
		videoInfo.Formats = append(videoInfo.Formats, format)
	}

	sortFormatsByResolution(videoInfo.Formats)
}

func buildQualityMap(formats []interface{}) map[string]models.VideoFormat {
	qualityMap := make(map[string]models.VideoFormat)

	for i, f := range formats {
		format, ok := f.(map[string]interface{})
		if !ok {
			continue
		}

		if i < 10 {
			logFormatDetails(i, format)
		}

		if shouldSkipFormat(format) {
			continue
		}

		videoFormat, resolution := buildVideoFormat(format)
		if _, exists := qualityMap[resolution]; !exists {
			qualityMap[resolution] = videoFormat
		}
	}

	return qualityMap
}

func logFormatDetails(index int, format map[string]interface{}) {
	log.Printf(consts.LOG_FORMAT_DETAILS,
		index, 
		format[consts.JSON_HEIGHT], 
		format[consts.JSON_VCODEC], 
		format[consts.JSON_ACODEC], 
		format[consts.JSON_EXT], 
		format[consts.JSON_FORMAT_ID], 
		format[consts.JSON_TBR])
}

func shouldSkipFormat(format map[string]interface{}) bool {
	vcodec, ok := format[consts.JSON_VCODEC].(string)
	return ok && vcodec == consts.VCODEC_NONE
}

func buildVideoFormat(format map[string]interface{}) (models.VideoFormat, string) {
	var videoFormat models.VideoFormat
	resolution := consts.RESOLUTION_UNKNOWN

	if height, ok := format[consts.JSON_HEIGHT].(float64); ok {
		resolution = fmt.Sprintf(consts.RESOLUTION_FORMAT, int(height))
		log.Printf(consts.LOG_FOUND_VIDEO_FORMAT, resolution, format[consts.JSON_FORMAT_ID])
	}

	videoFormat.Resolution = resolution

	if formatID, ok := format[consts.JSON_FORMAT_ID].(string); ok {
		videoFormat.FormatID = formatID
	}

	if ext, ok := format[consts.JSON_EXT].(string); ok {
		videoFormat.Extension = ext
	}

	setVideoFormatFileSize(&videoFormat, format)

	return videoFormat, resolution
}

func setVideoFormatFileSize(videoFormat *models.VideoFormat, format map[string]interface{}) {
	if filesize, ok := format[consts.JSON_FILESIZE].(float64); ok {
		if filesize > 0 {
			videoFormat.FileSize = formatFileSize(int64(filesize))
		} else if filesizeApprox, ok := format[consts.JSON_FILESIZE_APPROX].(float64); ok {
			videoFormat.FileSize = consts.APPROX_PREFIX + formatFileSize(int64(filesizeApprox))
		}
	}
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