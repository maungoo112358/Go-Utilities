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
    TIMEOUTS, 
    REGEX_PATTERNS 
} from './constants.js';

const API_BASE = API_ENDPOINTS.BASE;
let currentVideoInfo = null;
let currentDownloadId = null;

export function initVideoDownloader() {
    const urlInput = document.getElementById(ELEMENT_IDS.URL_INPUT);
    const resolutionSection = document.getElementById(ELEMENT_IDS.RESOLUTION_SECTION);
    const confirmDownloadBtn = document.getElementById(ELEMENT_IDS.CONFIRM_DOWNLOAD_BTN);
    
    console.log('Video downloader elements:', {
        urlInput: !!urlInput,
        resolutionSection: !!resolutionSection,
        confirmDownloadBtn: !!confirmDownloadBtn
    });
    
    let debounceTimer;
    if (urlInput) {
        urlInput.addEventListener('input', (e) => {
            clearTimeout(debounceTimer);
            const url = e.target.value.trim();
            
            if (url === '') {
                hideResolutionSection();
                currentVideoInfo = null; 
                return;
            }
            
            debounceTimer = setTimeout(() => {
                if (isValidYouTubeURL(url)) {
                    fetchVideoInfo(url);
                } else {
                    hideResolutionSection();
                }
            }, TIMEOUTS.DEBOUNCE_INPUT);
        });
    }
    
    if (confirmDownloadBtn) {
        confirmDownloadBtn.addEventListener('click', handleConfirmDownload);
    }
    
    document.addEventListener('click', (e) => {
        if (e.target.id === ELEMENT_IDS.CANCEL_BTN) {
            handleCancelDownload();
        } else if (e.target.id === ELEMENT_IDS.PAUSE_RESUME_BTN) {
            handlePauseResumeDownload();
        }
    });
}

export function isValidYouTubeURL(url) {
    return REGEX_PATTERNS.YOUTUBE_URL.test(url);
}

async function fetchVideoInfo(url) {
    const resolutionSection = document.getElementById(ELEMENT_IDS.RESOLUTION_SECTION);
    
    try {
        showLoadingState();
        
        const response = await fetch(`${API_BASE}${API_ENDPOINTS.VIDEO_INFO}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ url }),
        });
        
        const data = await response.json();
        
        if (response.ok) {
            currentVideoInfo = data;
            populateResolutions(data.formats);
            showResolutionSection();
        } else {
            window.showError(data.message || ERROR_MESSAGES.FAILED_FETCH_VIDEO_INFO);
            hideResolutionSection();
        }
    } catch (error) {
        window.showError(ERROR_MESSAGES.CONNECTION_ERROR);
        hideResolutionSection();
    }
}

function showLoadingState() {
    const resolutionSection = document.getElementById('resolutionSection');
    resolutionSection.innerHTML = `
        <div style="text-align: center; padding: 20px;">
            <span class="loading-spinner"></span>
            <span style="margin-left: 8px; color: #BBBBBB;">${UI_TEXT.FETCHING_VIDEO_INFO}</span>
        </div>
    `;
    resolutionSection.classList.remove(CSS_CLASSES.HIDDEN);
}

function populateResolutions(formats) {
    const resolutionSection = document.getElementById(ELEMENT_IDS.RESOLUTION_SECTION);
    resolutionSection.innerHTML = `
        <label for="resolutionSelect" class="resolution-label">${UI_TEXT.SELECT_RESOLUTION_LABEL}</label>
        <select id="resolutionSelect" class="resolution-select">
            <option value="">${UI_TEXT.SELECT_RESOLUTION}</option>
        </select>
        <button id="confirmDownloadBtn" class="confirm-download-btn">
            ${UI_TEXT.START_DOWNLOAD}
        </button>
    `;
    
    const newResolutionSelect = document.getElementById(ELEMENT_IDS.RESOLUTION_SELECT);
    const newConfirmDownloadBtn = document.getElementById(ELEMENT_IDS.CONFIRM_DOWNLOAD_BTN);
    
    newResolutionSelect.addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && newResolutionSelect.value) {
            handleConfirmDownload();
        }
    });
    newConfirmDownloadBtn.addEventListener('click', handleConfirmDownload);
    
    const sortedFormats = formats.sort((a, b) => {
        const resA = parseInt(a.resolution) || 0;
        const resB = parseInt(b.resolution) || 0;
        return resB - resA;
    });
    
    sortedFormats.forEach(format => {
        const option = document.createElement('option');
        option.value = format.format_id;
        option.textContent = `${format.resolution} ${format.ext ? `(${format.ext.toUpperCase()})` : ''}${format.filesize ? ` - ${format.filesize}` : ''}`;
        newResolutionSelect.appendChild(option);
    });
}

function showResolutionSection() {
    const resolutionSection = document.getElementById(ELEMENT_IDS.RESOLUTION_SECTION);
    resolutionSection.classList.remove(CSS_CLASSES.HIDDEN);
}

function hideResolutionSection() {
    const resolutionSection = document.getElementById(ELEMENT_IDS.RESOLUTION_SECTION);
    resolutionSection.classList.add(CSS_CLASSES.HIDDEN);
}

async function handleConfirmDownload() {
    const resolutionSelect = document.getElementById(ELEMENT_IDS.RESOLUTION_SELECT);
    const confirmDownloadBtn = document.getElementById(ELEMENT_IDS.CONFIRM_DOWNLOAD_BTN);
    const selectedQuality = resolutionSelect.value;
    
    if (!selectedQuality) {
        window.showError(ERROR_MESSAGES.SELECT_RESOLUTION);
        return;
    }
    
    if (!currentVideoInfo) {
        window.showError(ERROR_MESSAGES.VIDEO_INFO_NOT_AVAILABLE);
        return;
    }
    
    confirmDownloadBtn.innerHTML = `<span class="loading-spinner"></span>${UI_TEXT.STARTING}`;
    confirmDownloadBtn.disabled = true;
    
    try {
        const response = await fetch(`${API_BASE}${API_ENDPOINTS.DOWNLOAD}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ 
                url: currentVideoInfo.parsed_url, 
                quality: selectedQuality 
            }),
        });
        
        const data = await response.json();
        
        if (data.success) {
            currentDownloadId = data.filename;
            window.setCurrentDownloadId(data.filename);
            showProgress();
            hideResolutionSection();
        } else {
            window.showError(data.message || ERROR_MESSAGES.DOWNLOAD_FAILED);
        }
    } catch (error) {
        window.showError(ERROR_MESSAGES.CONNECTION_ERROR);
    } finally {
        confirmDownloadBtn.innerHTML = UI_TEXT.START_DOWNLOAD;
        confirmDownloadBtn.disabled = false;
    }
}

function showProgress() {
    const progressContainer = document.getElementById(ELEMENT_IDS.PROGRESS_CONTAINER);
    const urlInput = document.getElementById(ELEMENT_IDS.URL_INPUT);
    
    progressContainer.classList.remove(CSS_CLASSES.HIDDEN);
    urlInput.disabled = true;
    urlInput.style.opacity = '0.5';
}

export function hideProgress() {
    const progressContainer = document.getElementById(ELEMENT_IDS.PROGRESS_CONTAINER);
    const urlInput = document.getElementById(ELEMENT_IDS.URL_INPUT);
    
    progressContainer.classList.add(CSS_CLASSES.HIDDEN);
    currentDownloadId = null;
    
    urlInput.disabled = false;
    urlInput.style.opacity = '1';
    
    const progressFill = document.querySelector(SELECTORS.PROGRESS_FILL);
    const progressPercentage = document.querySelector(SELECTORS.PROGRESS_PERCENTAGE);
    const downloadSpeed = document.querySelector(SELECTORS.DOWNLOAD_SPEED);
    const downloadEta = document.querySelector(SELECTORS.DOWNLOAD_ETA);
    const progressText = document.querySelector(SELECTORS.PROGRESS_TEXT);
    
    if (progressFill) progressFill.style.width = '0%';
    if (progressPercentage) progressPercentage.textContent = '0%';
    if (downloadSpeed) downloadSpeed.textContent = UI_TEXT.SPEED_PLACEHOLDER;
    if (downloadEta) downloadEta.textContent = UI_TEXT.ETA_PLACEHOLDER;
    if (progressText) progressText.textContent = UI_TEXT.DOWNLOADING;
    
    if (currentVideoInfo && document.getElementById(ELEMENT_IDS.URL_INPUT).value.trim() && isValidYouTubeURL(document.getElementById(ELEMENT_IDS.URL_INPUT).value.trim())) {
        showResolutionSection();
    }
}

async function handleCancelDownload() {
    if (!currentDownloadId) return;
    
    try {
        const response = await fetch(`${API_BASE}${API_ENDPOINTS.CANCEL}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ downloadId: currentDownloadId }),
        });
        
        if (response.ok) {
            window.showError(ERROR_MESSAGES.DOWNLOAD_CANCELLED);
            hideProgress();
        }
    } catch (error) {
        console.error(LOG_MESSAGES.FAILED_CANCEL_DOWNLOAD, error);
        window.showError(ERROR_MESSAGES.FAILED_CANCEL_DOWNLOAD);
    }
}

let isPaused = false;

async function handlePauseResumeDownload() {
    if (!currentDownloadId) return;
    
    const pauseResumeBtn = document.getElementById(ELEMENT_IDS.PAUSE_RESUME_BTN);
    if (!pauseResumeBtn) return;
    
    try {
        const action = isPaused ? 'resume' : 'pause';
        const endpoint = isPaused ? API_ENDPOINTS.RESUME : API_ENDPOINTS.PAUSE;
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: HTTP_METHODS.POST,
            headers: {
                'Content-Type': CONTENT_TYPES.JSON,
            },
            body: JSON.stringify({ downloadId: currentDownloadId }),
        });
        
        if (response.ok) {
            isPaused = !isPaused;
            pauseResumeBtn.textContent = isPaused ? UI_TEXT.RESUME : UI_TEXT.PAUSE;
            pauseResumeBtn.className = isPaused ? 'control-btn pause-btn resume-btn' : 'control-btn pause-btn';
            
            const progressText = document.querySelector(SELECTORS.PROGRESS_TEXT);
            if (progressText) {
                progressText.textContent = isPaused ? UI_TEXT.PAUSED : UI_TEXT.DOWNLOADING;
            }
        }
    } catch (error) {
        console.error(LOG_MESSAGES.FAILED_PAUSE_RESUME, error);
        window.showError(ERROR_MESSAGES.FAILED_PAUSE_RESUME_DOWNLOAD);
    }
}

export function getCurrentVideoInfo() {
    return currentVideoInfo;
}