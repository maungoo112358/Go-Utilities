import { initVideoDownloader, hideProgress, isValidYouTubeURL, getCurrentVideoInfo } from './video_downloader.js';
import { initAudioConverter, hideMp3Progress, handleMp3ProgressUpdate } from './audio_converter.js';
import { initJsonFormatter } from './json_formatter.js';
import { 
    LOG_MESSAGES, 
    ERROR_MESSAGES, 
    UI_TEXT, 
    APP_TITLES, 
    CSS_CLASSES, 
    ELEMENT_IDS, 
    SELECTORS, 
    API_ENDPOINTS, 
    WS_MESSAGE_TYPES, 
    DOWNLOAD_STATUS, 
    TIMEOUTS, 
    SHUTDOWN_MESSAGES 
} from './constants.js';

const API_BASE = API_ENDPOINTS.BASE;
let ws = null;
let currentDownloadId = null;
let isPaused = false;

document.addEventListener('DOMContentLoaded', () => {
    const progressContainer = document.getElementById('progressContainer');
    
    initWebSocket();
    initMenuSystem();
    initVideoDownloader();
    console.log('About to call initAudioConverter...');
    initAudioConverter();
    console.log('initAudioConverter called');
    initJsonFormatter();
    
    function initMenuSystem() {
        const menuButtons = document.querySelectorAll('.menu-btn');
        const apps = document.querySelectorAll('.app');
        
        menuButtons.forEach(button => {
            button.addEventListener('click', () => {
                const targetApp = button.getAttribute('data-app');
                
                menuButtons.forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
                
                apps.forEach(app => {
                    if (app.classList.contains(`${targetApp}-app`)) {
                        app.classList.remove('hidden');
                    } else {
                        app.classList.add('hidden');
                    }
                });
                
                if (currentDownloadId && targetApp !== 'youtube-video') {
                    cancelCurrentDownload();
                }
                
                const titles = {
                    'youtube-video': APP_TITLES.YOUTUBE_VIDEO,
                    'youtube-mp3': APP_TITLES.YOUTUBE_MP3, 
                    'json-formatter': APP_TITLES.JSON_FORMATTER
                };
                document.title = titles[targetApp] || APP_TITLES.DEFAULT;
            });
        });
    }
    
    function initWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws`;
        
        ws = new WebSocket(wsUrl);
        
        ws.onopen = () => {
            console.log(LOG_MESSAGES.WS_CONNECTED);
        };
        
        ws.onmessage = (event) => {
            const update = JSON.parse(event.data);
            console.log(LOG_MESSAGES.WS_MESSAGE, update);
            
            if (update.type === WS_MESSAGE_TYPES.SHUTDOWN) {
                console.log(LOG_MESSAGES.WS_SHUTDOWN_SIGNAL);
                showShutdownMessage();
                setTimeout(() => {
                    window.close();
                    if (!window.closed) {
                        window.location.href = 'about:blank';
                    }
                }, TIMEOUTS.SHUTDOWN_CLOSE_TAB);
                return;
            }
            
            handleProgressUpdate(update);
        };
        
        ws.onerror = (error) => {
            console.error(LOG_MESSAGES.WS_ERROR, error);
        };
        
        ws.onclose = () => {
            console.log(LOG_MESSAGES.WS_DISCONNECTED);
            setTimeout(initWebSocket, TIMEOUTS.WS_RECONNECT);
        };
    }
    
    function handleProgressUpdate(update) {
        console.log(LOG_MESSAGES.HANDLING_PROGRESS_UPDATE, update.id, LOG_MESSAGES.CURRENT_DOWNLOAD, currentDownloadId);
        
        if (update.id !== currentDownloadId) return;
        
        const isMp3 = update.id.startsWith('mp3_');
        
        if (isMp3) {
            handleMp3ProgressUpdate(update);
        } else {
            handleVideoProgressUpdate(update);
        }
    }
    
    function handleVideoProgressUpdate(update) {
        const progressFill = document.querySelector(`#${ELEMENT_IDS.PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_FILL}`);
        const progressPercentage = document.querySelector(`#${ELEMENT_IDS.PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_PERCENTAGE}`);
        const progressText = document.querySelector(`#${ELEMENT_IDS.PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_TEXT}`);
        const downloadSpeed = document.querySelector(`#${ELEMENT_IDS.PROGRESS_CONTAINER} ${SELECTORS.DOWNLOAD_SPEED}`);
        const downloadEta = document.querySelector(`#${ELEMENT_IDS.PROGRESS_CONTAINER} ${SELECTORS.DOWNLOAD_ETA}`);
        
        if (progressFill) progressFill.style.width = `${update.progress}%`;
        if (progressPercentage) progressPercentage.textContent = `${Math.round(update.progress)}%`;
        
        if (update.speed && downloadSpeed) {
            downloadSpeed.textContent = update.speed;
        }
        
        if (update.eta && downloadEta) {
            downloadEta.textContent = `${UI_TEXT.ETA_PREFIX}${update.eta}`;
        }
        
        switch (update.status) {
            case DOWNLOAD_STATUS.DOWNLOADING:
                if (progressText) progressText.textContent = UI_TEXT.DOWNLOADING;
                break;
            case DOWNLOAD_STATUS.PROCESSING:
                if (progressText) progressText.textContent = UI_TEXT.PROCESSING;
                break;
            case DOWNLOAD_STATUS.COMPLETED:
                if (progressText) progressText.textContent = UI_TEXT.COMPLETED;
                setTimeout(() => {
                    hideProgress();
                }, TIMEOUTS.AUTO_HIDE_PROGRESS);
                break;
            case DOWNLOAD_STATUS.ERROR:
                showError(update.message || ERROR_MESSAGES.DOWNLOAD_FAILED);
                hideProgress();
                break;
        }
    }
    
    function cancelCurrentDownload() {
        if (!currentDownloadId) return;
        
        const isMp3 = currentDownloadId.startsWith('mp3_');
        
        if (isMp3) {
            document.getElementById(ELEMENT_IDS.MP3_CANCEL_BTN)?.click();
        } else {
            document.getElementById(ELEMENT_IDS.CANCEL_BTN)?.click();
        }
    }
    
    
    
    
    function showShutdownMessage() {
        const overlay = document.createElement('div');
        overlay.className = CSS_CLASSES.SHUTDOWN_OVERLAY;
        
        const message = document.createElement('div');
        message.className = CSS_CLASSES.SHUTDOWN_MESSAGE;
        
        const title = document.createElement('h2');
        title.className = CSS_CLASSES.SHUTDOWN_TITLE;
        title.textContent = SHUTDOWN_MESSAGES.TITLE;
        
        const text = document.createElement('p');
        text.className = CSS_CLASSES.SHUTDOWN_TEXT;
        text.textContent = SHUTDOWN_MESSAGES.MESSAGE;
        
        message.appendChild(title);
        message.appendChild(text);
        overlay.appendChild(message);
        document.body.appendChild(overlay);
    }
    
    window.showError = function(message) {
        const toast = document.createElement('div');
        toast.className = CSS_CLASSES.ERROR_TOAST;
        toast.textContent = message;
        
        document.body.appendChild(toast);
        
        setTimeout(() => {
            toast.style.animation = 'slideOut 0.3s ease';
            setTimeout(() => toast.remove(), TIMEOUTS.ERROR_TOAST_SLIDE_OUT);
        }, TIMEOUTS.ERROR_TOAST_DURATION);
    };
    
    window.setCurrentDownloadId = function(id) {
        currentDownloadId = id;
    };
    
    window.getCurrentDownloadId = function() {
        return currentDownloadId;
    };
});

