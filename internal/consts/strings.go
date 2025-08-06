package consts

// ---------- LOG MESSAGES - INFORMATIONAL --------------
const (
	LOG_YT_DLP_VERSION           = "yt-dlp version: %s"
	LOG_PATH                     = "Path: %s"
	LOG_FULL_ARGS                = "Full args: %v"
	LOG_QUALITY_REQUESTED        = "Quality requested: %s"
	LOG_URL                      = "URL: %s"
	LOG_URL_MP3                  = "URL: %s -> %s"
	LOG_TEMP_DIR_MP3             = "Temp dir: %s"
	LOG_DOWNLOADING_COMMAND      = "=== DOWNLOADING WITH COMMAND ==="
	LOG_CONVERTING_TO_MP3_PATH   = "=== CONVERTING TO MP3 === | Path: %s"
	LOG_RAW_VIDEO_INFO           = "=== RAW VIDEO INFO OUTPUT ==="
	LOG_AVAILABLE_FORMATS_TOTAL  = "=== AVAILABLE FORMATS === | Total formats found: %d"
	LOG_RAW_VIDEO_INFO_LENGTH    = "=== RAW VIDEO INFO OUTPUT === | Output length: %d bytes"
	LOG_FIRST_1000_CHARS        = "First 1000 chars: %s"
	LOG_FULL_OUTPUT              = "Full output: %s"
	LOG_YT_DLP_INFO_FAILED_STDERR = "yt-dlp video info command failed: %v | yt-dlp stderr: %s"
	LOG_FORMAT_DETAILS           = "Format %d: height=%v, vcodec=%v, acodec=%v, ext=%v, format_id=%v, tbr=%v"
	LOG_FOUND_VIDEO_FORMAT       = "Found video format: %s (format_id: %v)"
	LOG_GETTING_VIDEO_INFO       = "Getting video info with command: %s -j --no-warnings %s"
	LOG_BROADCASTING_UPDATE      = "Broadcasting update: %+v"
	LOG_BROADCASTING_SUBSCRIBERS = "Broadcasting to %d subscribers"
	LOG_SENT_UPDATE              = "Sent update to subscriber %d"
	LOG_SUBSCRIBER_CHANNEL_FULL  = "Subscriber %d channel full, skipping"
	LOG_SERVER_STARTING          = "Server starting on %s"
	LOG_UNSUPPORTED_PLATFORM     = "Unsupported platform, please open %s manually"
	LOG_OPEN_MANUALLY            = "Please open %s manually"
	LOG_OPENING_BROWSER          = "Opening %s in your default browser..."
	LOG_RECEIVED_SIGNAL          = "Received signal: %v. Shutting down gracefully..."
	LOG_SHUTDOWN_COMPLETE        = "Server shutdown complete"
	LOG_OPENING_FILE_EXPLORER    = "Opening File Explorer..."
)

// ---------- LOG MESSAGES - WARNINGS --------------
const (
	WARNING_YT_DLP_OUTDATED      = "WARNING: yt-dlp version may be outdated. Consider updating from https://github.com/yt-dlp/yt-dlp/releases"
	WARNING_FFMPEG_NOT_FOUND     = "Warning: FFmpeg not found, audio merging may not work: %v"
	WARNING_FFMPEG_NOT_FOUND_MP3 = "Warning: FFmpeg not found, MP3 conversion may not work: %v"
	LOG_YT_DLP_TEST_FAILED       = "WARNING: yt-dlp test failed: %v"
)

// ---------- LOG MESSAGES - WEBSOCKET --------------
const (
	LOG_WS_UPGRADE_ERROR          = "WebSocket upgrade error: %v"
	LOG_WS_CONNECTION_ESTABLISHED = "WebSocket connection established"
	LOG_SENDING_WS_UPDATE         = "Sending WebSocket update: %+v"
	LOG_WS_WRITE_ERROR            = "WebSocket write error: %v"
	LOG_SENDING_SHUTDOWN_TO_WS    = "Sending shutdown signal to WebSocket client"
	LOG_SENDING_SHUTDOWN_SIGNAL   = "Sending shutdown signal to all WebSocket clients..."
	LOG_SHUTDOWN_SIGNAL_SENT      = "Shutdown signal sent"
	LOG_SHUTDOWN_SIGNAL_FULL      = "Shutdown signal channel full"
)

// ---------- LOG MESSAGES - PROCESS OUTPUT --------------
const (
	LOG_YT_DLP_STDERR         = "yt-dlp stderr: %s"
	LOG_YT_DLP_MP3_STDERR     = "yt-dlp MP3 stderr: %s"
	LOG_YT_DLP_STDOUT         = "yt-dlp stdout: %s"
	LOG_CMD_WAIT_FAILED_OUTPUT = "cmd.Wait() failed with error: %v | Full stdout output: %v"
	LOG_EXIT_CODE_STDERR      = "Exit code: %d, stderr: %s"
	LOG_EXIT_CODE_101_SUCCESS = "yt-dlp exit code 101 due to --max-downloads or existing file, treating as success"
	LOG_PROCESS_COMPLETED_LOOKING = "yt-dlp process completed successfully | Looking for converted file in: %s"
	LOG_FILES_FOUND_TEMP_DIR  = "Files found in temp dir: %v"
)

// ---------- LOG MESSAGES - DOWNLOAD/CONVERSION PROCESS --------------
const (
	LOG_DOWNLOAD_COMMAND_INFO    = "=== DOWNLOADING WITH COMMAND === | Path: %s | Full args: %v | Quality requested: %s | URL: %s"
	LOG_STARTING_DOWNLOAD        = "Starting download for URL: %s, Quality: %s"
	LOG_DOWNLOAD_STARTED         = "Download started with ID: %s"
	LOG_STARTING_MP3_CONVERSION  = "Starting MP3 conversion for URL: %s"
	LOG_MP3_CONVERSION_STARTED   = "MP3 conversion started with ID: %s"
	LOG_INVALID_REQUEST_BODY     = "Invalid request body: %v"
	LOG_INVALID_REQUEST_BODY_MP3 = "Invalid request body: %v"
	LOG_TEMPLATE_ERROR           = "Template error: %v"
	LOG_TEMPLATE_EXECUTION_ERROR = "Template execution error: %v"
)

// ---------- PROGRESS/STATUS MESSAGES --------------
const (
	MSG_STARTING_DOWNLOAD       = "Starting download..."
	MSG_DOWNLOAD_COMPLETE       = "Download completed, processing..."
	MSG_CONVERTING_VIDEO        = "Converting video..."
	MSG_STARTING_MP3_CONVERSION = "Starting MP3 conversion..."
	MSG_MP3_CONVERSION_COMPLETE = "MP3 conversion completed, processing..."
	MSG_CONVERTING_TO_MP3       = "Converting to MP3..."
	MSG_SAVED_AS                = "Saved as: %s"
	MSG_MP3_SAVED_AS            = "MP3 saved as: %s"
)

// ---------- USER NOTIFICATION MESSAGES --------------
const (
	DUPLICATE_MSG         = "Duplicate file detected. Saving as: %s"
	POWERSHELL_FAILED     = "PowerShell SaveFileDialog failed: %v"
	MP3_CONVERSION_FAILED = "MP3 conversion failed: %v"
)

// ---------- ERROR MESSAGES - FILE SYSTEM --------------
const (
	ERR_YT_DLP_NOT_FOUND     = "yt-dlp.exe not found at %s"
	ERR_FFMPEG_NOT_FOUND     = "ffmpeg.exe not found at %s"
	ERR_NO_VIDEO_FILE        = "no video file found in %s"
	ERR_CREATE_TEMP_DIR      = "Failed to create temp directory: %v"
	ERR_RENAME_FILE          = "Failed to rename file with resolution: %v"
	ERR_SAVE_FILE            = "Failed to save file: %v"
	ERR_SAVE_MP3_FILE        = "Failed to save MP3 file: %v"
	ERR_SAVE_FILE_PICKER     = "failed to save file: %v"
	ERR_SAVE_CANCELLED       = "save cancelled by user"
	ERR_FIND_DOWNLOADED_FILE = "Could not find downloaded file: %v"
	ERR_FIND_MP3_FILE        = "Could not find converted MP3 file: %v"
)

// ---------- ERROR MESSAGES - PROCESS EXECUTION --------------
const (
	ERR_CREATE_STDOUT_PIPE     = "Failed to create stdout pipe: %v"
	ERR_CREATE_STDERR_PIPE     = "Failed to create stderr pipe: %v"
	ERR_CREATE_STDOUT_PIPE_MP3 = "Failed to create stdout pipe: %v"
	ERR_CREATE_STDERR_PIPE_MP3 = "Failed to create stderr pipe: %v"
	ERR_START_YT_DLP           = "Failed to start yt-dlp: %v"
	ERR_START_YT_DLP_EXE       = "Failed to start yt-dlp: %v. Make sure yt-dlp.exe exists and is executable"
	ERR_START_YT_DLP_DOWNLOAD  = "yt-dlp not found: %v"
	ERR_START_YT_DLP_MP3       = "Failed to start yt-dlp for MP3: %v"
	ERR_START_YT_DLP_INFO      = "yt-dlp not found: %v"
	ERR_START_MP3_CONVERSION   = "Failed to start MP3 conversion: %v"
	ERR_YT_DLP_INFO_FAILED     = "yt-dlp video info command failed: %v"
	ERR_YT_DLP_TEST_FAILED     = "yt-dlp test failed: %v"
)

// ---------- ERROR MESSAGES - URL AND VIDEO HANDLING --------------
const (
	ERR_INVALID_URL          = "invalid URL: %v"
	ERR_NOT_YOUTUBE_URL      = "not a YouTube URL"
	ERR_EXTRACT_VIDEO_ID     = "could not extract video ID from URL"
	ERR_MISSING_V_PARAM      = "missing v parameter in URL"
	ERR_INVALID_YOUTUBE_URL  = "invalid YouTube URL: %v"
	ERR_VIDEO_UNAVAILABLE    = "video is unavailable or private"
	ERR_YOUTUBE_BLOCKED      = "youtube blocked the request. Try a different video or wait a moment."
	ERR_VIDEO_RESTRICTED     = "video requires sign-in or is age-restricted"
	ERR_VIDEO_RESTRICTED_GEO = "video is restricted or geo-blocked. YouTube blocked access to this video"
	ERR_VIDEO_FRAGMENTS      = "video fragments are unavailable. This video may be corrupted or restricted"
	ERR_VIDEO_PRIVATE        = "this is a private video"
	ERR_VIDEO_REMOVED        = "video is unavailable or has been removed"
	ERR_DOWNLOAD_FAILED      = "Download failed: %s"
	ERR_PARSE_VIDEO_INFO     = "failed to parse video info: %v"
	ERR_GET_VIDEO_INFO       = "failed to get video info: %s"
	ERR_GET_VIDEO_INFO_2     = "failed to get video info: %v"
)

// ---------- ERROR MESSAGES - HTTP AND REQUESTS --------------
const (
	ERR_INVALID_REQUEST      = "Invalid request"
	ERR_INVALID_REQUEST_INFO = "Invalid request"
	ERR_INVALID_REQUEST_MP3  = "Invalid request"
	ERR_TEMPLATE             = "Template error: %s"
	ERR_TEMPLATE_EXECUTION   = "Template execution error"
	ERR_SERVER_START         = "Server failed to start:"
	ERR_OPEN_BROWSER         = "Failed to open browser automatically: %v"
	LOG_BROWSER_OPEN_FAILED  = "Failed to open browser automatically: %v | Please open %s manually"
	ERR_FORCED_SHUTDOWN      = "Server forced to shutdown: %v"
	ERR_SEND_SHUTDOWN_SIGNAL = "Failed to send shutdown signal: %v"
)

// ---------- YT-DLP OUTPUT TEXT PATTERNS --------------
const (
	YT_DLP_DOWNLOAD_100_PERCENT     = "[download] 100%"
	YT_DLP_ALREADY_DOWNLOADED       = "has already been downloaded"
	YT_DLP_FFMPEG_TAG               = "[ffmpeg]"
	YT_DLP_MAX_DOWNLOADS_REACHED    = "Maximum number of downloads reached"
	YT_DLP_VIDEO_UNAVAILABLE        = "Video unavailable"
	YT_DLP_FORBIDDEN_403            = "403"
	YT_DLP_FORBIDDEN_TEXT           = "Forbidden"
	YT_DLP_SIGN_IN_REQUIRED         = "Sign in"
	YT_DLP_FRAGMENT_TEXT            = "fragment"
	YT_DLP_NOT_FOUND_TEXT           = "not found"
	YT_DLP_PRIVATE_VIDEO_TEXT       = "Private video"
	MP3_CONVERSION_FAILED_PREFIX    = "mp3 conversion failed: %s"
	MP3_CONVERSION_FAILED_EXIT_CODE = "mp3 conversion failed (exit code %d): %s"
	MP3_CONVERSION_FAILED_GENERIC   = "mp3 conversion failed: %v"
)

//---------- URL TEMPLATES --------------
const (
	YOUTUBE_WATCH_URL = "https://www.youtube.com/watch?v=%s"
)

//---------- HTTP RESPONSE MESSAGES --------------
const (
	MSG_DOWNLOAD_STARTED     = "Download started"
	MSG_MP3_CONVERSION_STARTED = "MP3 conversion started"
	MSG_SHUTDOWN_SIGNAL      = "Application is shutting down"
	MSG_TAB_CLOSE_AUTO       = "This tab will close automatically."
	MSG_APP_SHUTTING_DOWN    = "Application Shutting Down"
)


//nolint:ST1005
