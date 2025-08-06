package consts

//---------- SERVER CONFIGURATION --------------
const (
	DEFAULT_PORT       = ":8484"
	BASE_URL           = "http://localhost"
	FILE_PICKER_FILTER = "Video Files|*.mp4;*.mkv;*.avi;*.mov;*.wmv;*.flv;*.webm|All Files|*.*"
	SAVE_DIALOG_TITLE  = "Save Video As"
)

//---------- EXECUTABLE NAMES AND FILE EXTENSIONS --------------
const (
	YT_DLP_EXE_NAME  = "yt-dlp.exe"
	FFMPEG_EXE_NAME  = "ffmpeg.exe"
	FRAGMENT_EXT     = ".f"
	TEMP_EXT         = ".temp"
	DEPENDENCIES_DIR = "dependencies"
)

//---------- TEMPORARY DIRECTORY NAMES --------------
const (
	TEMP_DIR = "go-utilities-temp"
)

//---------- YT-DLP CONFIGURATION --------------
const (
	YT_DLP_OUTPUT_FORMAT = "%(title)s.%(ext)s"
	USER_AGENT_STRING    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
)

//---------- HTTP HEADERS FOR YT-DLP --------------
const (
	HEADER_ACCEPT          = "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	HEADER_ACCEPT_LANGUAGE = "Accept-Language: en-us,en;q=0.5"
	HEADER_ACCEPT_ENCODING = "Accept-Encoding: gzip,deflate"
	HEADER_ACCEPT_CHARSET  = "Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7"
	HEADER_CONNECTION      = "Connection: keep-alive"
)

//---------- USER INTERFACE MESSAGES --------------
const (
	CANCELLED_MESSAGE = "CANCELLED"
	SHUTDOWN_TITLE    = "Shutting Down"
)

//---------- APPLICATION STATE CONSTANTS --------------
const (
	STATUS_STARTING    = "starting"
	STATUS_DOWNLOADING = "downloading"
	STATUS_CONVERTING  = "converting"
	STATUS_ERROR       = "error"
	STATUS_COMPLETED   = "completed"
)

//---------- FORMAT AND ID TEMPLATES --------------
const (
	DOWNLOAD_ID_FORMAT = "dl_%d"
	MP3_ID_FORMAT      = "mp3_%d"
	TIMESTAMP_FORMAT   = "20060102-150405"
)

//---------- APPLICATION DEFAULTS --------------
const (
	BEST_QUALITY = "best"
)

//---------- SYSTEM COMMANDS --------------
const (
	POWERSHELL_COMMAND = "powershell"
	COMMAND_FLAG       = "-Command"
	EXPLORER_COMMAND   = "explorer"
	WINDOWS_OS         = "windows"
)

//---------- YT-DLP COMMAND OPTIONS --------------
const (
	FFMPEG_LOCATION_FLAG = "--ffmpeg-location"
	YT_DLP_VERSION_FLAG  = "--version"
)

//---------- URL PATTERNS AND COMPONENTS --------------
const (
	YOUTUBE_DOMAIN     = "youtube.com"
	YOUTU_BE_DOMAIN    = "youtu.be"
	YOUTUBE_EMBED_PATH = "/embed/"
	URL_PATH_SEPARATOR = "/"
	URL_VIDEO_PARAM    = "v"
)

//---------- VERSION AND VALIDATION --------------
const (
	MIN_YTDLP_YEAR = "2024"
)

//---------- FILE SIZE FORMATTING --------------
const (
	BYTES_UNIT      = 1024
	FILE_SIZE_UNITS = "KMGTPE"
)

//---------- POWERSHELL SCRIPTS --------------
const (
	POWERSHELL_FILE_PICKER_SCRIPT = `
Add-Type -AssemblyName System.Windows.Forms
$saveFileDialog = New-Object System.Windows.Forms.SaveFileDialog
$saveFileDialog.Filter = "%s"
$saveFileDialog.FileName = "%s"
$saveFileDialog.Title = "%s"
$saveFileDialog.DefaultExt = "%s"

if ($saveFileDialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $saveFileDialog.FileName
} else {
    Write-Output "%s"
}`
)

//---------- YT-DLP OUTPUT PARSING PATTERNS --------------
const (
	YT_DLP_PROGRESS_REGEX_WITH_SPEED = `\[download\]\s+(\d+\.?\d*)%\s+of\s+.*?\s+at\s+(\S+)\s+ETA\s+(\S+)`
	YT_DLP_PROGRESS_REGEX_SIMPLE     = `\[download\]\s+(\d+\.?\d*)%`
	YT_DLP_TITLE_REGEX               = `\[download\] Destination: (.+)`
)