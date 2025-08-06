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
	FORMAT_FLAG          = "-f"
)

//---------- YT-DLP FORMAT STRINGS --------------
const (
	QUALITY_SUFFIX                 = "p"
	QUALITY_HEIGHT_FORMAT          = "bestvideo[height<=%s][ext=mp4]+bestaudio[ext=m4a]/bestvideo[height<=%s]+bestaudio/best[height<=%s]"
	QUALITY_CUSTOM_FORMAT          = "(%s+bestaudio[ext=m4a])/(%s+bestaudio)/%s/best"
	QUALITY_BEST_FORMAT            = "bestvideo[height>=1080]+bestaudio[ext=m4a]/bestvideo[height>=1080]+bestaudio/bestvideo[height>=720][fps>=30]+bestaudio[ext=m4a]/bestvideo[height>=720][fps>=30]+bestaudio/bestvideo[height>=720]+bestaudio[ext=m4a]/bestvideo[height>=720]+bestaudio/best[height>=720]/best"
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
	APPROX_PREFIX   = "~"
)

//---------- VIDEO INFO JSON FIELDS --------------
const (
	JSON_TITLE            = "title"
	JSON_DURATION         = "duration"
	JSON_THUMBNAIL        = "thumbnail"
	JSON_FORMATS          = "formats"
	JSON_HEIGHT           = "height"
	JSON_VCODEC           = "vcodec"
	JSON_ACODEC           = "acodec"
	JSON_EXT              = "ext"
	JSON_FORMAT_ID        = "format_id"
	JSON_TBR              = "tbr"
	JSON_FILESIZE         = "filesize"
	JSON_FILESIZE_APPROX  = "filesize_approx"
	VCODEC_NONE           = "none"
	RESOLUTION_UNKNOWN    = "unknown"
	RESOLUTION_FORMAT     = "%dp"
	DURATION_FORMAT       = "%d:%02d"
)

//---------- HTTP ROUTES AND PATHS --------------
const (
	STATIC_DIR_PATH           = "./static"
	STATIC_ROUTE_PREFIX       = "/static/"
	HOME_ROUTE                = "/"
	SHUTDOWN_ROUTE            = "/shutdown"
	API_ROUTE_PREFIX          = "/api"
	DOWNLOAD_ROUTE            = "/download"
	MP3_CONVERT_ROUTE         = "/mp3-convert"
	VIDEO_INFO_ROUTE          = "/video-info"
	WEBSOCKET_ROUTE           = "/ws"
	TEMPLATE_PATH             = "static/html/index.html"
	SHUTDOWN_TEMPLATE_PATH    = "static/html/shutdown.html"
)

//---------- HTTP METHODS --------------
const (
	HTTP_GET  = "GET"
	HTTP_POST = "POST"
)

//---------- HTTP HEADERS --------------
const (
	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_HTML = "text/html"
	HEADER_CONTENT_TYPE = "Content-Type"
)

//---------- WEBSOCKET CONFIGURATION --------------
const (
	SHUTDOWN_SIGNAL_BUFFER = 10
	SHUTDOWN_DELAY_MS      = 500
	WS_MESSAGE_TYPE_SHUTDOWN = "shutdown"
)

//---------- BROWSER PATHS --------------
const (
	CHROME_PATH_1    = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	CHROME_PATH_2    = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
	CHROME_PATH_USER = "\\Google\\Chrome\\Application\\chrome.exe"
	
	FIREFOX_PATH_1 = "C:\\Program Files\\Mozilla Firefox\\firefox.exe"
	FIREFOX_PATH_2 = "C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe"
	
	EDGE_PATH_1 = "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	EDGE_PATH_2 = "C:\\Program Files\\Microsoft\\Edge\\Application\\msedge.exe"
	
	BRAVE_PATH_1 = "C:\\Program Files\\BraveSoftware\\Brave-Browser\\Application\\brave.exe"
	BRAVE_PATH_2 = "C:\\Program Files (x86)\\BraveSoftware\\Brave-Browser\\Application\\brave.exe"
	
	BROWSER_NAME_CHROME  = "chrome"
	BROWSER_NAME_FIREFOX = "firefox"
	BROWSER_NAME_EDGE    = "edge"
	BROWSER_NAME_BRAVE   = "brave"
	
	RUNDLL32_COMMAND = "rundll32"
	URL_DLL_HANDLER  = "url.dll,FileProtocolHandler"
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