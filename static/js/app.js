const API_BASE = '/api';
let ws = null;
let currentDownloadId = null;
let currentVideoInfo = null;
let isPaused = false;

document.addEventListener('DOMContentLoaded', () => {
    const urlInput = document.getElementById('urlInput');
    const resolutionSection = document.getElementById('resolutionSection');
    const resolutionSelect = document.getElementById('resolutionSelect');
    const confirmDownloadBtn = document.getElementById('confirmDownloadBtn');
    const progressContainer = document.getElementById('progressContainer');
    
    // MP3 elements
    const mp3UrlInput = document.getElementById('mp3UrlInput');
    const convertMp3Btn = document.getElementById('convertMp3Btn');
    const mp3ProgressContainer = document.getElementById('mp3ProgressContainer');
    
    // Initialize WebSocket
    initWebSocket();
    
    // Initialize menu system
    initMenuSystem();
    
    // URL input change handler with debounce
    let debounceTimer;
    urlInput.addEventListener('input', (e) => {
        clearTimeout(debounceTimer);
        const url = e.target.value.trim();
        
        if (url === '') {
            hideResolutionSection();
            currentVideoInfo = null; // Clear video info when URL is cleared
            return;
        }
        
        // Debounce to avoid too many API calls
        debounceTimer = setTimeout(() => {
            if (isValidYouTubeURL(url)) {
                fetchVideoInfo(url);
            } else {
                hideResolutionSection();
            }
        }, 800); // Wait 800ms after user stops typing
    });
    
    // Confirm download button click handler
    confirmDownloadBtn.addEventListener('click', handleConfirmDownload);
    
    // Enter key handler for resolution select
    resolutionSelect.addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && resolutionSelect.value) {
            handleConfirmDownload();
        }
    });
    
    // Control button handlers (using event delegation since buttons are created dynamically)
    document.addEventListener('click', (e) => {
        if (e.target.id === 'cancelBtn') {
            handleCancelDownload();
        } else if (e.target.id === 'pauseResumeBtn') {
            handlePauseResumeDownload();
        } else if (e.target.id === 'convertMp3Btn') {
            handleMp3Convert();
        } else if (e.target.id === 'mp3CancelBtn') {
            handleCancelDownload(); // Reuse the same cancel logic
        } else if (e.target.id === 'mp3PauseResumeBtn') {
            handlePauseResumeDownload(); // Reuse the same pause/resume logic
        }
    });
    
    function initMenuSystem() {
        const menuButtons = document.querySelectorAll('.menu-btn');
        const apps = document.querySelectorAll('.app');
        
        menuButtons.forEach(button => {
            button.addEventListener('click', () => {
                const targetApp = button.getAttribute('data-app');
                
                // Update active menu button
                menuButtons.forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
                
                // Show target app, hide others
                apps.forEach(app => {
                    if (app.classList.contains(`${targetApp}-app`)) {
                        app.classList.remove('hidden');
                    } else {
                        app.classList.add('hidden');
                    }
                });
                
                // Reset any ongoing downloads when switching apps
                if (currentDownloadId && targetApp !== 'youtube-video') {
                    handleCancelDownload();
                }
                
                // Update page title
                const titles = {
                    'youtube-video': 'YouTube Video Downloader',
                    'youtube-mp3': 'YouTube Video to MP3 Downloader', 
                    'json-formatter': 'JSON Formatter'
                };
                document.title = titles[targetApp] || 'Go Utilities';
            });
        });
    }
    
    function isValidYouTubeURL(url) {
        const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube\.com\/(watch\?v=|embed\/|v\/)|youtu\.be\/)[\w-]+/;
        return youtubeRegex.test(url);
    }
    
    async function fetchVideoInfo(url) {
        try {
            // Show loading in resolution section
            showLoadingState();
            
            const response = await fetch(`${API_BASE}/video-info`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url }),
            });
            
            const data = await response.json();
            
            if (response.ok) {
                currentVideoInfo = data;
                populateResolutions(data.formats);
                showResolutionSection();
            } else {
                showError(data.message || 'Failed to fetch video information');
                hideResolutionSection();
            }
        } catch (error) {
            showError('Connection error. Please try again.');
            hideResolutionSection();
        }
    }
    
    function showLoadingState() {
        resolutionSection.innerHTML = `
            <div style="text-align: center; padding: 20px;">
                <span class="loading-spinner"></span>
                <span style="margin-left: 8px; color: #BBBBBB;">Fetching video information...</span>
            </div>
        `;
        resolutionSection.classList.remove('hidden');
    }
    
    function populateResolutions(formats) {
        resolutionSection.innerHTML = `
            <label for="resolutionSelect" class="resolution-label">Select Resolution:</label>
            <select id="resolutionSelect" class="resolution-select">
                <option value="">Select resolution...</option>
            </select>
            <button id="confirmDownloadBtn" class="confirm-download-btn">
                START DOWNLOAD
            </button>
        `;
        
        const newResolutionSelect = document.getElementById('resolutionSelect');
        const newConfirmDownloadBtn = document.getElementById('confirmDownloadBtn');
        
        // Re-attach event listeners
        newResolutionSelect.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && newResolutionSelect.value) {
                handleConfirmDownload();
            }
        });
        newConfirmDownloadBtn.addEventListener('click', handleConfirmDownload);
        
        // Sort formats by resolution (highest first)
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
        resolutionSection.classList.remove('hidden');
    }
    
    function hideResolutionSection() {
        resolutionSection.classList.add('hidden');
        // Don't clear currentVideoInfo - keep it for after download
    }
    
    async function handleConfirmDownload() {
        const resolutionSelect = document.getElementById('resolutionSelect');
        const confirmDownloadBtn = document.getElementById('confirmDownloadBtn');
        const selectedQuality = resolutionSelect.value;
        
        if (!selectedQuality) {
            showError('Please select a resolution');
            return;
        }
        
        if (!currentVideoInfo) {
            showError('Video information not available');
            return;
        }
        
        // Show loading state
        confirmDownloadBtn.innerHTML = '<span class="loading-spinner"></span>STARTING...';
        confirmDownloadBtn.disabled = true;
        
        try {
            const response = await fetch(`${API_BASE}/download`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    url: currentVideoInfo.parsed_url, 
                    quality: selectedQuality 
                }),
            });
            
            const data = await response.json();
            
            if (data.success) {
                currentDownloadId = data.filename;
                showProgress();
                hideResolutionSection(); // Hide during download
                // Keep URL and video info for after download
            } else {
                showError(data.message || 'Download failed');
            }
        } catch (error) {
            showError('Connection error. Please try again.');
        } finally {
            // Reset button
            confirmDownloadBtn.innerHTML = 'START DOWNLOAD';
            confirmDownloadBtn.disabled = false;
        }
    }
    
    function initWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws`;
        
        ws = new WebSocket(wsUrl);
        
        ws.onopen = () => {
            console.log('WebSocket connected');
        };
        
        ws.onmessage = (event) => {
            const update = JSON.parse(event.data);
            console.log('Progress update:', update); // Debug logging
            handleProgressUpdate(update);
        };
        
        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        
        ws.onclose = () => {
            console.log('WebSocket disconnected');
            setTimeout(initWebSocket, 3000);
        };
    }
    
    function handleProgressUpdate(update) {
        console.log('Handling progress update for:', update.id, 'Current download:', currentDownloadId);
        
        if (update.id !== currentDownloadId) return;
        
        const progressFill = document.querySelector('.progress-fill');
        const progressPercentage = document.querySelector('.progress-percentage');
        const progressText = document.querySelector('.progress-text');
        const downloadSpeed = document.querySelector('.download-speed');
        const downloadEta = document.querySelector('.download-eta');
        
        if (progressFill) progressFill.style.width = `${update.progress}%`;
        if (progressPercentage) progressPercentage.textContent = `${Math.round(update.progress)}%`;
        
        if (update.speed && downloadSpeed) {
            downloadSpeed.textContent = update.speed;
        }
        
        if (update.eta && downloadEta) {
            downloadEta.textContent = `ETA: ${update.eta}`;
        }
        
        switch (update.status) {
            case 'downloading':
                if (progressText) progressText.textContent = 'Downloading...';
                break;
            case 'processing':
                if (progressText) progressText.textContent = 'Processing...';
                break;
            case 'completed':
                if (progressText) progressText.textContent = 'Completed!';
                setTimeout(() => {
                    hideProgress();
                    // Show resolution section again if URL is still valid
                    if (currentVideoInfo && urlInput.value.trim() && isValidYouTubeURL(urlInput.value.trim())) {
                        showResolutionSection();
                    }
                }, 3000);
                break;
            case 'error':
                showError(update.message || 'Download failed');
                hideProgress();
                break;
        }
    }
    
    function showProgress() {
        progressContainer.classList.remove('hidden');
        // Disable URL input during download
        urlInput.disabled = true;
        urlInput.style.opacity = '0.5';
    }
    
    function hideProgress() {
        progressContainer.classList.add('hidden');
        currentDownloadId = null;
        
        // Re-enable URL input
        urlInput.disabled = false;
        urlInput.style.opacity = '1';
        
        // Reset progress
        const progressFill = document.querySelector('.progress-fill');
        const progressPercentage = document.querySelector('.progress-percentage');
        const downloadSpeed = document.querySelector('.download-speed');
        const downloadEta = document.querySelector('.download-eta');
        const progressText = document.querySelector('.progress-text');
        
        if (progressFill) progressFill.style.width = '0%';
        if (progressPercentage) progressPercentage.textContent = '0%';
        if (downloadSpeed) downloadSpeed.textContent = '0 MB/s';
        if (downloadEta) downloadEta.textContent = 'ETA: --:--';
        if (progressText) progressText.textContent = 'Downloading...';
    }
    
    async function handleCancelDownload() {
        if (!currentDownloadId) return;
        
        try {
            const response = await fetch(`${API_BASE}/cancel`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ downloadId: currentDownloadId }),
            });
            
            if (response.ok) {
                showError('Download cancelled');
                hideProgress();
            }
        } catch (error) {
            console.error('Failed to cancel download:', error);
            showError('Failed to cancel download');
        }
    }
    
    async function handlePauseResumeDownload() {
        if (!currentDownloadId) return;
        
        const pauseResumeBtn = document.getElementById('pauseResumeBtn');
        if (!pauseResumeBtn) return;
        
        try {
            const action = isPaused ? 'resume' : 'pause';
            const response = await fetch(`${API_BASE}/${action}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ downloadId: currentDownloadId }),
            });
            
            if (response.ok) {
                isPaused = !isPaused;
                pauseResumeBtn.textContent = isPaused ? 'RESUME' : 'PAUSE';
                pauseResumeBtn.className = isPaused ? 'control-btn pause-btn resume-btn' : 'control-btn pause-btn';
                
                const progressText = document.querySelector('.progress-text');
                if (progressText) {
                    progressText.textContent = isPaused ? 'Paused' : 'Downloading...';
                }
            }
        } catch (error) {
            console.error('Failed to pause/resume download:', error);
            showError('Failed to pause/resume download');
        }
    }
    
    async function handleMp3Convert() {
        const mp3UrlInput = document.getElementById('mp3UrlInput');
        const convertMp3Btn = document.getElementById('convertMp3Btn');
        const mp3ProgressContainer = document.getElementById('mp3ProgressContainer');
        
        if (!mp3UrlInput || !mp3UrlInput.value.trim()) {
            showError('Please enter a YouTube URL');
            return;
        }
        
        if (!isValidYouTubeURL(mp3UrlInput.value.trim())) {
            showError('Please enter a valid YouTube URL');
            return;
        }
        
        // Show loading state
        convertMp3Btn.innerHTML = '<span class="loading-spinner"></span>STARTING...';
        convertMp3Btn.disabled = true;
        
        try {
            const response = await fetch(`${API_BASE}/mp3-convert`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    url: mp3UrlInput.value.trim()
                }),
            });
            
            const data = await response.json();
            
            if (data.success) {
                currentDownloadId = data.filename;
                showMp3Progress();
                // Disable URL input during conversion
                mp3UrlInput.disabled = true;
                mp3UrlInput.style.opacity = '0.5';
            } else {
                showError(data.message || 'MP3 conversion failed');
            }
        } catch (error) {
            showError('Connection error. Please try again.');
        } finally {
            // Reset button
            convertMp3Btn.innerHTML = 'CONVERT TO MP3';
            convertMp3Btn.disabled = false;
        }
    }
    
    function showMp3Progress() {
        const mp3ProgressContainer = document.getElementById('mp3ProgressContainer');
        if (mp3ProgressContainer) {
            mp3ProgressContainer.classList.remove('hidden');
        }
    }
    
    function hideMp3Progress() {
        const mp3ProgressContainer = document.getElementById('mp3ProgressContainer');
        const mp3UrlInput = document.getElementById('mp3UrlInput');
        
        if (mp3ProgressContainer) {
            mp3ProgressContainer.classList.add('hidden');
        }
        
        // Re-enable URL input
        if (mp3UrlInput) {
            mp3UrlInput.disabled = false;
            mp3UrlInput.style.opacity = '1';
        }
        
        currentDownloadId = null;
    }
    
    function showError(message) {
        const toast = document.createElement('div');
        toast.className = 'error-toast';
        toast.textContent = message;
        
        document.body.appendChild(toast);
        
        setTimeout(() => {
            toast.style.animation = 'slideOut 0.3s ease';
            setTimeout(() => toast.remove(), 300);
        }, 4000);
    }
});