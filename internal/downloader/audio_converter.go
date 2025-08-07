package downloader

import (
	"Go-Utilities/internal/consts"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ExecuteMp3Conversion(url string, progressCallback ProgressCallback) (*YtDlpResult, error) {
	cleanURL, tempDir, err := prepareMp3ConversionEnvironment(url)
	if err != nil {
		return nil, err
	}

	args, err := buildMp3ConversionCommand(tempDir, cleanURL)
	if err != nil {
		return nil, err
	}

	title, stdoutLines, err := executeMp3ConversionProcess(args, progressCallback)
	if err != nil {
		validationErr := validateMp3ConversionResult(err, stdoutLines)
		if validationErr != nil {
			return nil, validationErr
		}
	}

	return locateMp3ConversionResult(tempDir, title)
}

func prepareMp3ConversionEnvironment(url string) (string, string, error) {
	cleanURL, err := cleanYouTubeURL(url)
	if err != nil {
		return "", "", fmt.Errorf(consts.ERR_INVALID_YOUTUBE_URL, err)
	}

	tempDir := filepath.Join(os.TempDir(), consts.TEMP_DIR)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", "", fmt.Errorf(consts.ERR_CREATE_TEMP_DIR, err)
	}

	return cleanURL, tempDir, nil
}

func buildMp3ConversionCommand(tempDir, cleanURL string) ([]string, error) {
	outputPath := filepath.Join(tempDir, consts.YT_DLP_OUTPUT_FORMAT)

	ffmpegPath, err := getFFmpegPath()
	if err != nil {
		log.Printf(consts.WARNING_FFMPEG_NOT_FOUND_MP3, err)
	}

	args := []string{"-o", outputPath}
	args = append(args, consts.YT_DLP_MP3_ARGS...)
	args = append(args, cleanURL)

	if ffmpegPath != "" {
		args = append(args, consts.FFMPEG_LOCATION_FLAG, ffmpegPath)
	}

	return args, nil
}

func executeMp3ConversionProcess(args []string, progressCallback ProgressCallback) (string, []string, error) {
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		return "", nil, fmt.Errorf(consts.ERR_START_YT_DLP_MP3, err)
	}

	log.Printf(consts.LOG_CONVERTING_TO_MP3_PATH, ytDlpPath)

	cmd := exec.Command(ytDlpPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", nil, fmt.Errorf(consts.ERR_CREATE_STDOUT_PIPE_MP3, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", nil, fmt.Errorf(consts.ERR_CREATE_STDERR_PIPE_MP3, err)
	}

	if err := cmd.Start(); err != nil {
		return "", nil, fmt.Errorf(consts.ERR_START_MP3_CONVERSION, err)
	}

	go handleStderrOutput(stderr)

	scanner := bufio.NewScanner(stdout)
	title, stdoutLines := handleMp3ProcessOutput(scanner, progressCallback)

	err = cmd.Wait()
	return title, stdoutLines, err
}

func handleStderrOutput(stderr io.ReadCloser) {
	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		line := stderrScanner.Text()
		log.Printf(consts.LOG_YT_DLP_MP3_STDERR, line)
	}
}

func handleMp3ProcessOutput(scanner *bufio.Scanner, progressCallback ProgressCallback) (string, []string) {
	progressRegex1 := regexp.MustCompile(consts.YT_DLP_PROGRESS_REGEX_WITH_SPEED)
	progressRegex2 := regexp.MustCompile(consts.YT_DLP_PROGRESS_REGEX_SIMPLE)
	titleRegex := regexp.MustCompile(consts.YT_DLP_TITLE_REGEX)

	var title, filename string
	var stdoutLines []string

	for scanner.Scan() {
		line := scanner.Text()
		stdoutLines = append(stdoutLines, line)
		log.Printf(consts.LOG_YT_DLP_STDOUT, line)

		if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
			filename = filepath.Base(matches[1])
			title = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		parseMp3Progress(line, progressRegex1, progressRegex2, progressCallback)
	}

	return title, stdoutLines
}

func parseMp3Progress(line string, progressRegex1, progressRegex2 *regexp.Regexp, progressCallback ProgressCallback) {
	if progressCallback == nil {
		return
	}

	if matches := progressRegex1.FindStringSubmatch(line); len(matches) > 3 {
		progress := parseFloat(matches[1])
		speed := matches[2]
		eta := matches[3]
		progressCallback(progress, speed, eta, consts.MSG_STARTING_MP3_CONVERSION)
	} else if matches := progressRegex2.FindStringSubmatch(line); len(matches) > 1 {
		progress := parseFloat(matches[1])
		progressCallback(progress, "", "", consts.MSG_STARTING_MP3_CONVERSION)
	}

	if strings.Contains(line, consts.YT_DLP_DOWNLOAD_100_PERCENT) || strings.Contains(line, consts.YT_DLP_ALREADY_DOWNLOADED) {
		progressCallback(100, "", "", consts.MSG_MP3_CONVERSION_COMPLETE)
	}

	if strings.Contains(line, consts.YT_DLP_FFMPEG_TAG) {
		progressCallback(100, "", "", consts.MSG_CONVERTING_TO_MP3)
	}
}

func validateMp3ConversionResult(err error, stdoutLines []string) error {
	log.Printf(consts.LOG_CMD_WAIT_FAILED_OUTPUT, err, stdoutLines)

	if exitError, ok := err.(*exec.ExitError); ok {
		stderrOutput := string(exitError.Stderr)
		log.Printf(consts.LOG_EXIT_CODE_STDERR, exitError.ExitCode(), stderrOutput)

		fullOutput := strings.Join(stdoutLines, "\n")

		if exitError.ExitCode() == 101 && (strings.Contains(fullOutput, consts.YT_DLP_MAX_DOWNLOADS_REACHED) || strings.Contains(fullOutput, consts.YT_DLP_ALREADY_DOWNLOADED)) {
			log.Printf(consts.LOG_EXIT_CODE_101_SUCCESS)
			return nil
		}

		if strings.Contains(fullOutput, consts.YT_DLP_VIDEO_UNAVAILABLE) || strings.Contains(stderrOutput, consts.YT_DLP_VIDEO_UNAVAILABLE) {
			return fmt.Errorf(consts.ERR_VIDEO_UNAVAILABLE)
		}

		if strings.Contains(fullOutput, consts.YT_DLP_FORBIDDEN_403) || strings.Contains(fullOutput, consts.YT_DLP_FORBIDDEN_TEXT) || strings.Contains(stderrOutput, consts.YT_DLP_FORBIDDEN_403) || strings.Contains(stderrOutput, consts.YT_DLP_FORBIDDEN_TEXT) {
			return fmt.Errorf(consts.ERR_YOUTUBE_BLOCKED)
		}

		if strings.Contains(fullOutput, consts.YT_DLP_SIGN_IN_REQUIRED) || strings.Contains(stderrOutput, consts.YT_DLP_SIGN_IN_REQUIRED) {
			return fmt.Errorf(consts.ERR_VIDEO_RESTRICTED)
		}

		if stderrOutput != "" {
			return fmt.Errorf(consts.MP3_CONVERSION_FAILED_PREFIX, stderrOutput)
		}

		return fmt.Errorf(consts.MP3_CONVERSION_FAILED_EXIT_CODE, exitError.ExitCode(), strings.Join(stdoutLines[len(stdoutLines)-5:], "; "))
	}

	return fmt.Errorf(consts.MP3_CONVERSION_FAILED_GENERIC, err)
}

func locateMp3ConversionResult(tempDir, title string) (*YtDlpResult, error) {
	log.Printf(consts.LOG_PROCESS_COMPLETED_LOOKING, tempDir)

	if files, err := filepath.Glob(filepath.Join(tempDir, "*")); err == nil {
		log.Printf(consts.LOG_FILES_FOUND_TEMP_DIR, files)
	}

	convertedFile, err := findDownloadedFile(tempDir)
	if err != nil {
		return nil, fmt.Errorf(consts.ERR_FIND_MP3_FILE, err)
	}

	return &YtDlpResult{
		Title:    title,
		FilePath: convertedFile,
		TempDir:  tempDir,
		Success:  true,
	}, nil
}
