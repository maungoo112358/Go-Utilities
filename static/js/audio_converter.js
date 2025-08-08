import { 
    LOG_MESSAGES, 
    ERROR_MESSAGES, 
    UI_TEXT, 
    CSS_CLASSES, 
    ELEMENT_IDS, 
    SELECTORS, 
    API_ENDPOINTS, 
    CONTENT_TYPES, 
    HTTP_METHODS, 
    DOWNLOAD_STATUS, 
    TIMEOUTS, 
    REGEX_PATTERNS 
} from './constants.js';

const API_BASE = API_ENDPOINTS.BASE;

export function initAudioConverter() {
    console.log('initAudioConverter called');
    const mp3UrlInput = document.getElementById(ELEMENT_IDS.MP3_URL_INPUT);
    const convertMp3Btn = document.getElementById(ELEMENT_IDS.CONVERT_MP3_BTN);
    
    console.log('MP3 elements found:', {
        mp3UrlInput: !!mp3UrlInput,
        convertMp3Btn: !!convertMp3Btn,
        mp3UrlInputId: mp3UrlInput?.id,
        convertMp3BtnId: convertMp3Btn?.id
    });
    
    if (!mp3UrlInput || !convertMp3Btn) {
        console.error(LOG_MESSAGES.MP3_CONVERTER_ELEMENTS_NOT_FOUND);
        return;
    }
    
    convertMp3Btn.addEventListener('click', function(e) {
        console.log('MP3 button clicked! Event:', e);
        console.log(LOG_MESSAGES.MP3_BUTTON_CLICKED_DIRECT);
        e.preventDefault();
        e.stopPropagation();
        handleMp3Convert();
    });
    
    console.log('MP3 button event listener added');
    
    mp3UrlInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && mp3UrlInput.value.trim()) {
            handleMp3Convert();
        }
    });
    
    document.addEventListener('click', (e) => {
        if (e.target.id === ELEMENT_IDS.MP3_CANCEL_BTN) {
            console.log(LOG_MESSAGES.MP3_CANCEL_CLICKED);
            handleCancelMp3Download();
        } else if (e.target.id === ELEMENT_IDS.MP3_PAUSE_RESUME_BTN) {
            console.log(LOG_MESSAGES.MP3_PAUSE_RESUME_CLICKED);
            handlePauseResumeMp3Download();
        }
    });
}

export function isValidYouTubeURL(url) {
    return REGEX_PATTERNS.YOUTUBE_URL.test(url);
}

async function handleMp3Convert() {
    console.log(LOG_MESSAGES.MP3_CONVERTER_CALLED);
    
    const mp3UrlInput = document.getElementById(ELEMENT_IDS.MP3_URL_INPUT);
    const convertMp3Btn = document.getElementById(ELEMENT_IDS.CONVERT_MP3_BTN);
    const mp3ProgressContainer = document.getElementById(ELEMENT_IDS.MP3_PROGRESS_CONTAINER);
    
    console.log(LOG_MESSAGES.MP3_ELEMENTS_FOUND, {
        urlInput: !!mp3UrlInput,
        button: !!convertMp3Btn,
        progressContainer: !!mp3ProgressContainer,
        urlValue: mp3UrlInput?.value
    });
    
    if (!mp3UrlInput || !mp3UrlInput.value.trim()) {
        console.log(LOG_MESSAGES.NO_URL_ENTERED);
        window.showError(ERROR_MESSAGES.ENTER_YOUTUBE_URL);
        return;
    }
    
    if (!isValidYouTubeURL(mp3UrlInput.value.trim())) {
        window.showError(ERROR_MESSAGES.ENTER_VALID_YOUTUBE_URL);
        return;
    }
    
    convertMp3Btn.innerHTML = `<span class="loading-spinner"></span>${UI_TEXT.STARTING}`;
    convertMp3Btn.disabled = true;
    
    try {
        const response = await fetch(`${API_BASE}${API_ENDPOINTS.MP3_CONVERT}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ 
                url: mp3UrlInput.value.trim()
            }),
        });
        
        const data = await response.json();
        
        if (data.success) {
            window.setCurrentDownloadId(data.filename);
            showMp3Progress();
            mp3UrlInput.disabled = true;
            mp3UrlInput.style.opacity = '0.5';
        } else {
            window.showError(data.message || ERROR_MESSAGES.MP3_CONVERSION_FAILED);
        }
    } catch (error) {
        window.showError(ERROR_MESSAGES.CONNECTION_ERROR);
    } finally {
        convertMp3Btn.innerHTML = UI_TEXT.CONVERT_TO_MP3;
        convertMp3Btn.disabled = false;
    }
}

function showMp3Progress() {
    const mp3ProgressContainer = document.getElementById(ELEMENT_IDS.MP3_PROGRESS_CONTAINER);
    if (mp3ProgressContainer) {
        mp3ProgressContainer.classList.remove(CSS_CLASSES.HIDDEN);
    }
}

export function hideMp3Progress() {
    const mp3ProgressContainer = document.getElementById(ELEMENT_IDS.MP3_PROGRESS_CONTAINER);
    const mp3UrlInput = document.getElementById(ELEMENT_IDS.MP3_URL_INPUT);
    
    if (mp3ProgressContainer) {
        mp3ProgressContainer.classList.add(CSS_CLASSES.HIDDEN);
    }
    
    if (mp3UrlInput) {
        mp3UrlInput.disabled = false;
        mp3UrlInput.style.opacity = '1';
    }
    
    window.setCurrentDownloadId(null);
}

async function handleCancelMp3Download() {
    const currentDownloadId = window.getCurrentDownloadId();
    if (!currentDownloadId || !currentDownloadId.startsWith('mp3_')) return;
    
    try {
        const response = await fetch(`${API_BASE}${API_ENDPOINTS.CANCEL}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ downloadId: currentDownloadId }),
        });
        
        if (response.ok) {
            window.showError(ERROR_MESSAGES.MP3_CONVERSION_CANCELLED);
            hideMp3Progress();
        }
    } catch (error) {
        console.error(LOG_MESSAGES.FAILED_CANCEL_MP3, error);
        window.showError(ERROR_MESSAGES.FAILED_CANCEL_MP3_CONVERSION);
    }
}

let isMp3Paused = false;

async function handlePauseResumeMp3Download() {
    const currentDownloadId = window.getCurrentDownloadId();
    if (!currentDownloadId || !currentDownloadId.startsWith('mp3_')) return;
    
    const pauseResumeBtn = document.getElementById(ELEMENT_IDS.MP3_PAUSE_RESUME_BTN);
    if (!pauseResumeBtn) return;
    
    try {
        const action = isMp3Paused ? 'resume' : 'pause';
        const endpoint = isMp3Paused ? API_ENDPOINTS.RESUME : API_ENDPOINTS.PAUSE;
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ downloadId: currentDownloadId }),
        });
        
        if (response.ok) {
            isMp3Paused = !isMp3Paused;
            pauseResumeBtn.textContent = isMp3Paused ? UI_TEXT.RESUME : UI_TEXT.PAUSE;
            pauseResumeBtn.className = isMp3Paused ? 'control-btn pause-btn resume-btn' : 'control-btn pause-btn';
            
            const progressText = document.querySelector(`#${ELEMENT_IDS.MP3_PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_TEXT}`);
            if (progressText) {
                progressText.textContent = isMp3Paused ? UI_TEXT.PAUSED : UI_TEXT.CONVERTING;
            }
        }
    } catch (error) {
        console.error(LOG_MESSAGES.FAILED_PAUSE_RESUME_MP3, error);
        window.showError(ERROR_MESSAGES.FAILED_PAUSE_RESUME_MP3);
    }
}

export function handleMp3ProgressUpdate(update) {
    const progressFill = document.querySelector(`#${ELEMENT_IDS.MP3_PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_FILL}`);
    const progressPercentage = document.querySelector(`#${ELEMENT_IDS.MP3_PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_PERCENTAGE}`);
    const progressText = document.querySelector(`#${ELEMENT_IDS.MP3_PROGRESS_CONTAINER} ${SELECTORS.PROGRESS_TEXT}`);
    const downloadSpeed = document.querySelector(`#${ELEMENT_IDS.MP3_PROGRESS_CONTAINER} ${SELECTORS.DOWNLOAD_SPEED}`);
    const downloadEta = document.querySelector(`#${ELEMENT_IDS.MP3_PROGRESS_CONTAINER} ${SELECTORS.DOWNLOAD_ETA}`);
    
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
        case DOWNLOAD_STATUS.CONVERTING:
            if (progressText) progressText.textContent = UI_TEXT.CONVERTING_TO_MP3;
            break;
        case DOWNLOAD_STATUS.PROCESSING:
            if (progressText) progressText.textContent = UI_TEXT.PROCESSING_MP3;
            break;
        case DOWNLOAD_STATUS.COMPLETED:
            if (progressText) progressText.textContent = UI_TEXT.COMPLETED;
            setTimeout(() => {
                hideMp3Progress();
            }, TIMEOUTS.AUTO_HIDE_PROGRESS);
            break;
        case DOWNLOAD_STATUS.ERROR:
            window.showError(update.message || ERROR_MESSAGES.MP3_CONVERSION_FAILED);
            hideMp3Progress();
            break;
    }
}