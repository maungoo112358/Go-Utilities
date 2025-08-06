// ---------- CONSOLE LOG MESSAGES --------------
export const LOG_MESSAGES = {
    WS_CONNECTED: 'WebSocket connected',
    WS_DISCONNECTED: 'WebSocket disconnected',
    WS_ERROR: 'WebSocket error:',
    WS_MESSAGE: 'WebSocket message:',
    WS_SHUTDOWN_SIGNAL: 'Received shutdown signal - closing tab',
    
    MP3_BUTTON_CHECK: 'MP3 button check:',
    MP3_BUTTON_DIRECT_CLICK: 'Direct MP3 button clicked!',
    MP3_BUTTON_EVENT_ADDED: 'Direct event listener added to MP3 button',
    MP3_BUTTON_CLICKED_DELEGATION: 'MP3 Convert button clicked via event delegation!',
    MP3_BUTTON_CLICKED_DIRECT: 'MP3 button clicked via direct listener',
    MP3_CANCEL_CLICKED: 'MP3 Cancel button clicked',
    MP3_PAUSE_RESUME_CLICKED: 'MP3 Pause/Resume button clicked',
    
    CANCEL_BUTTON_CLICKED: 'Cancel button clicked',
    PAUSE_RESUME_CLICKED: 'Pause/Resume button clicked',
    
    CLICK_DETECTED: 'Click detected:',
    UNHANDLED_CLICK: 'Unhandled click on:',
    
    MP3_CONVERTER_CALLED: 'handleMp3Convert function called',
    MP3_ELEMENTS_FOUND: 'MP3 elements found:',
    NO_URL_ENTERED: 'No URL entered',
    
    HANDLING_PROGRESS_UPDATE: 'Handling progress update for:',
    CURRENT_DOWNLOAD: 'Current download:',
    BROADCASTING_UPDATE: 'Broadcasting update:',
    BROADCASTING_SUBSCRIBERS: 'Broadcasting to subscribers',
    SENT_UPDATE: 'Sent update to subscriber',
    SUBSCRIBER_CHANNEL_FULL: 'Subscriber channel full, skipping',
    SENDING_WS_UPDATE: 'Sending WebSocket update:',
    SENDING_SHUTDOWN_TO_WS: 'Sending shutdown signal to WebSocket client',
    SENDING_SHUTDOWN_SIGNAL: 'Sending shutdown signal to all WebSocket clients...',
    
    FAILED_CANCEL_DOWNLOAD: 'Failed to cancel download:',
    FAILED_PAUSE_RESUME: 'Failed to pause/resume download:',
    FAILED_CANCEL_MP3: 'Failed to cancel MP3 conversion:',
    FAILED_PAUSE_RESUME_MP3: 'Failed to pause/resume MP3 conversion:',
    
    MP3_CONVERTER_ELEMENTS_NOT_FOUND: 'MP3 converter elements not found',
    JSON_FORMATTER_ELEMENTS_NOT_FOUND: 'JSON formatter elements not found'
};

// ---------- ERROR MESSAGES --------------
export const ERROR_MESSAGES = {
    ENTER_YOUTUBE_URL: 'Please enter a YouTube URL',
    ENTER_VALID_YOUTUBE_URL: 'Please enter a valid YouTube URL',
    SELECT_RESOLUTION: 'Please select a resolution',
    VIDEO_INFO_NOT_AVAILABLE: 'Video information not available',
    CONNECTION_ERROR: 'Connection error. Please try again.',
    DOWNLOAD_FAILED: 'Download failed',
    MP3_CONVERSION_FAILED: 'MP3 conversion failed',
    DOWNLOAD_CANCELLED: 'Download cancelled',
    MP3_CONVERSION_CANCELLED: 'MP3 conversion cancelled',
    FAILED_CANCEL_DOWNLOAD: 'Failed to cancel download',
    FAILED_PAUSE_RESUME_DOWNLOAD: 'Failed to pause/resume download',
    FAILED_CANCEL_MP3_CONVERSION: 'Failed to cancel MP3 conversion',
    FAILED_PAUSE_RESUME_MP3: 'Failed to pause/resume MP3 conversion',
    FAILED_FETCH_VIDEO_INFO: 'Failed to fetch video information',
    NO_INPUT_TO_COPY: 'No input to copy',
    NO_OUTPUT_TO_COPY: 'No output to copy',
    NO_JSON_CONTENT_TO_DOWNLOAD: 'No JSON content to download',
    INVALID_JSON_CHECK_SYNTAX: 'Invalid JSON - please check your syntax'
};

// ---------- SUCCESS MESSAGES --------------
export const SUCCESS_MESSAGES = {
    INPUT_COPIED: 'Input copied!',
    OUTPUT_COPIED: 'Output copied!',
    JSON_DOWNLOADED: 'JSON downloaded!'
};

// ---------- UI TEXT CONSTANTS --------------
export const UI_TEXT = {
    STARTING: 'STARTING...',
    START_DOWNLOAD: 'START DOWNLOAD',
    CONVERT_TO_MP3: 'CONVERT TO MP3',
    PAUSE: 'PAUSE',
    RESUME: 'RESUME',
    CANCEL: 'CANCEL',
    
    DOWNLOADING: 'Downloading...',
    CONVERTING_TO_MP3: 'Converting to MP3...',
    PROCESSING_MP3: 'Processing MP3...',
    PROCESSING: 'Processing...',
    CONVERTING: 'Converting...',
    COMPLETED: 'Completed!',
    PAUSED: 'Paused',
    
    FETCHING_VIDEO_INFO: 'Fetching video information...',
    SELECT_RESOLUTION: 'Select resolution...',
    SELECT_RESOLUTION_LABEL: 'Select Resolution:',
    
    FORMATTED_JSON_PLACEHOLDER: 'Formatted JSON will appear here...',
    READY_TO_FORMAT: 'Ready to format',
    WAITING_FOR_INPUT: 'Waiting for input',
    VALID_JSON: '✅ Valid JSON',
    FORMATTED_SUCCESSFULLY: '✨ Formatted successfully',
    INVALID_JSON: '❌ Invalid JSON',
    INVALID_JSON_PREFIX: '❌ Invalid JSON: ',
    FIX_INPUT_TO_SEE_OUTPUT: 'Fix input to see output',
    
    CHARACTERS_SUFFIX: ' characters',
    ZERO_CHARACTERS: '0 characters',
    
    ETA_PREFIX: 'ETA: ',
    ETA_PLACEHOLDER: 'ETA: --:--',
    SPEED_PLACEHOLDER: '0 MB/s'
};

// ---------- APP TITLES --------------
export const APP_TITLES = {
    YOUTUBE_VIDEO: 'YouTube Video Downloader',
    YOUTUBE_MP3: 'YouTube Video to MP3 Downloader',
    JSON_FORMATTER: 'JSON Formatter',
    DEFAULT: 'Go Utilities'
};

// ---------- CSS CLASSES --------------
export const CSS_CLASSES = {
    HIDDEN: 'hidden',
    ERROR_TOAST: 'error-toast',
    CONTROL_BTN: 'control-btn',
    PAUSE_BTN: 'pause-btn',
    RESUME_BTN: 'resume-btn',
    JSON_BRACKET: 'json-bracket',
    JSON_KEY: 'json-key',
    JSON_BOOLEAN: 'json-boolean',
    JSON_NULL: 'json-null',
    JSON_NUMBER: 'json-number',
    JSON_STRING: 'json-string',
    JSON_COMMA: 'json-comma',
    ACTIVE: 'active',
    JSON_NOTIFICATION: 'json-notification',
    NOTIFICATION_SHOW: 'show',
    NOTIFICATION_HIDE: 'hide',
    SHUTDOWN_OVERLAY: 'shutdown-overlay',
    SHUTDOWN_MESSAGE: 'shutdown-message',
    SHUTDOWN_TITLE: 'shutdown-title',
    SHUTDOWN_TEXT: 'shutdown-text'
};

// ---------- HTML ELEMENT IDS --------------
export const ELEMENT_IDS = {
    URL_INPUT: 'urlInput',
    MP3_URL_INPUT: 'mp3UrlInput',
    RESOLUTION_SECTION: 'resolutionSection',
    RESOLUTION_SELECT: 'resolutionSelect',
    CONFIRM_DOWNLOAD_BTN: 'confirmDownloadBtn',
    CONVERT_MP3_BTN: 'convertMp3Btn',
    PROGRESS_CONTAINER: 'progressContainer',
    MP3_PROGRESS_CONTAINER: 'mp3ProgressContainer',
    CANCEL_BTN: 'cancelBtn',
    PAUSE_RESUME_BTN: 'pauseResumeBtn',
    MP3_CANCEL_BTN: 'mp3CancelBtn',
    MP3_PAUSE_RESUME_BTN: 'mp3PauseResumeBtn',
    JSON_INPUT: 'jsonInput',
    JSON_OUTPUT: 'jsonOutput',
    JSON_ERROR: 'jsonError',
    INPUT_STATS: 'inputStats',
    OUTPUT_STATS: 'outputStats',
    INPUT_CHARS: 'inputChars',
    OUTPUT_CHARS: 'outputChars'
};

// ---------- CSS SELECTORS --------------
export const SELECTORS = {
    PROGRESS_FILL: '.progress-fill',
    PROGRESS_PERCENTAGE: '.progress-percentage',
    PROGRESS_TEXT: '.progress-text',
    DOWNLOAD_SPEED: '.download-speed',
    DOWNLOAD_ETA: '.download-eta',
    MENU_BTN: '.menu-btn',
    APP: '.app'
};

// ---------- API ENDPOINTS --------------
export const API_ENDPOINTS = {
    BASE: '/api',
    VIDEO_INFO: '/video-info',
    DOWNLOAD: '/download',
    MP3_CONVERT: '/mp3-convert',
    CANCEL: '/cancel',
    PAUSE: '/pause',
    RESUME: '/resume',
    WEBSOCKET: '/ws'
};

// ---------- HTTP METHODS --------------
export const HTTP_METHODS = {
    GET: 'GET',
    POST: 'POST'
};

// ---------- CONTENT TYPES --------------
export const CONTENT_TYPES = {
    JSON: 'application/json'
};

// ---------- DOWNLOAD STATUS --------------
export const DOWNLOAD_STATUS = {
    DOWNLOADING: 'downloading',
    CONVERTING: 'converting',
    PROCESSING: 'processing',
    COMPLETED: 'completed',
    ERROR: 'error'
};

// ---------- WEBSOCKET MESSAGE TYPES --------------
export const WS_MESSAGE_TYPES = {
    SHUTDOWN: 'shutdown'
};

// ---------- TIMEOUTS AND DELAYS --------------
export const TIMEOUTS = {
    DEBOUNCE_INPUT: 800,
    JSON_DEBOUNCE: 300,
    AUTO_HIDE_PROGRESS: 3000,
    SHUTDOWN_CLOSE_TAB: 2000,
    ERROR_TOAST_DURATION: 4000,
    ERROR_TOAST_SLIDE_OUT: 300,
    NOTIFICATION_SHOW: 100,
    NOTIFICATION_DURATION: 2000,
    NOTIFICATION_HIDE: 300,
    WS_RECONNECT: 3000
};

// ---------- SHUTDOWN MESSAGES --------------
export const SHUTDOWN_MESSAGES = {
    TITLE: 'Application Shutting Down',
    MESSAGE: 'This tab will close automatically...'
};

// ---------- FILE DOWNLOAD --------------
export const DOWNLOAD_CONFIG = {
    FILENAME_PREFIX: 'formatted_',
    FILENAME_EXTENSION: '.json',
    BLOB_TYPE: 'application/json'
};

// ---------- HTML COMPONENTS --------------
export const HTML_COMPONENTS = {
    JSON_BOOLEAN: (value) => `<span class="${CSS_CLASSES.JSON_BOOLEAN}">${value}</span>`,
    JSON_NULL: (value) => `<span class="${CSS_CLASSES.JSON_NULL}">${value}</span>`,
    JSON_NUMBER: (value) => `<span class="${CSS_CLASSES.JSON_NUMBER}">${value}</span>`,
    JSON_STRING: (value) => `<span class="${CSS_CLASSES.JSON_STRING}">${value}</span>`,
    JSON_BRACKET: (value) => `<span class="${CSS_CLASSES.JSON_BRACKET}">${value}</span>`,
    JSON_KEY: (value) => `<span class="${CSS_CLASSES.JSON_KEY}">${value}</span>`,
    JSON_COMMA: () => `<span class="${CSS_CLASSES.JSON_COMMA}">,</span>`
};

// ---------- REGEX PATTERNS --------------
export const REGEX_PATTERNS = {
    YOUTUBE_URL: /^(https?:\/\/)?(www\.)?(youtube\.com\/(watch\?v=|embed\/|v\/)|youtu\.be\/)[\w-]+/,
    NUMERIC_VALUE: /^-?\d+(\.\d+)?([eE][+-]?\d+)?$/,
    LEADING_SPACES: /^(\s*)/,
    JSON_KEY_QUOTES: /^"([^"]*)"$/
};