import { 
    LOG_MESSAGES, 
    ERROR_MESSAGES, 
    SUCCESS_MESSAGES, 
    UI_TEXT, 
    CSS_CLASSES, 
    ELEMENT_IDS, 
    DOWNLOAD_CONFIG, 
    TIMEOUTS, 
    REGEX_PATTERNS,
    HTML_COMPONENTS 
} from './constants.js';

let jsonDebounceTimer;

export function initJsonFormatter() {
    const jsonInput = document.getElementById(ELEMENT_IDS.JSON_INPUT);
    const jsonOutput = document.getElementById(ELEMENT_IDS.JSON_OUTPUT);
    
    if (!jsonInput || !jsonOutput) {
        console.error(LOG_MESSAGES.JSON_FORMATTER_ELEMENTS_NOT_FOUND);
        return;
    }
    
    jsonInput.addEventListener('input', beautifyJSON);
    
    window.beautifyJSON = beautifyJSON;
    window.copyJsonInput = copyJsonInput;
    window.copyJsonOutput = copyJsonOutput;
    window.downloadJSON = downloadJSON;
}

function beautifyJSON() {
    clearTimeout(jsonDebounceTimer);
    jsonDebounceTimer = setTimeout(() => {
        const input = document.getElementById(ELEMENT_IDS.JSON_INPUT).value.trim();
        const output = document.getElementById(ELEMENT_IDS.JSON_OUTPUT);
        const errorDiv = document.getElementById(ELEMENT_IDS.JSON_ERROR);
        const inputStats = document.getElementById(ELEMENT_IDS.INPUT_STATS);
        const outputStats = document.getElementById(ELEMENT_IDS.OUTPUT_STATS);
        const inputChars = document.getElementById(ELEMENT_IDS.INPUT_CHARS);
        const outputChars = document.getElementById(ELEMENT_IDS.OUTPUT_CHARS);
        
        if (inputChars) inputChars.textContent = input.length + UI_TEXT.CHARACTERS_SUFFIX;
        
        if (!input.trim()) {
            if (output) output.textContent = UI_TEXT.FORMATTED_JSON_PLACEHOLDER;
            if (errorDiv) errorDiv.classList.add(CSS_CLASSES.HIDDEN);
            if (inputStats) inputStats.textContent = UI_TEXT.READY_TO_FORMAT;
            if (outputStats) outputStats.textContent = UI_TEXT.WAITING_FOR_INPUT;
            if (outputChars) outputChars.textContent = UI_TEXT.ZERO_CHARACTERS;
            return;
        }
        
        try {
            const parsed = JSON.parse(input);
            const beautified = JSON.stringify(parsed, null, 2);
            
            const highlighted = highlightJSON(beautified);
            if (output) output.innerHTML = highlighted;
            if (errorDiv) errorDiv.classList.add(CSS_CLASSES.HIDDEN);
            
            if (inputStats) inputStats.textContent = UI_TEXT.VALID_JSON;
            if (outputStats) outputStats.textContent = UI_TEXT.FORMATTED_SUCCESSFULLY;
            if (outputChars) outputChars.textContent = beautified.length + UI_TEXT.CHARACTERS_SUFFIX;
            
        } catch (error) {
            if (output) output.textContent = ERROR_MESSAGES.INVALID_JSON_CHECK_SYNTAX;
            if (errorDiv) {
                errorDiv.textContent = UI_TEXT.INVALID_JSON_PREFIX + error.message;
                errorDiv.classList.remove(CSS_CLASSES.HIDDEN);
            }
            
            if (inputStats) inputStats.textContent = UI_TEXT.INVALID_JSON;
            if (outputStats) outputStats.textContent = UI_TEXT.FIX_INPUT_TO_SEE_OUTPUT;
            if (outputChars) outputChars.textContent = UI_TEXT.ZERO_CHARACTERS;
        }
    }, TIMEOUTS.JSON_DEBOUNCE);
}

function highlightJSON(json) {
    let result = json.replace(/&/g, '&amp;')
                    .replace(/</g, '&lt;')
                    .replace(/>/g, '&gt;');
    
    let lines = result.split('\n');
    let highlightedLines = [];
    
    for (let line of lines) {
        let highlightedLine = highlightJSONLine(line);
        highlightedLines.push(highlightedLine);
    }
    
    return highlightedLines.join('\n');
}

function highlightJSONLine(line) {
    if (line.trim() === '') return line;
    
    let leadingSpaces = line.match(/^(\s*)/)[0];
    let content = line.trim();
    
    if (content === '{' || content === '}' || content === '[' || content === ']' || 
        content === '{,' || content === '},' || content === '[,' || content === '],') {
        return leadingSpaces + HTML_COMPONENTS.JSON_BRACKET(content);
    }
    
    let colonIndex = content.indexOf(':');
    if (colonIndex > 0) {
        let keyPart = content.substring(0, colonIndex).trim();
        let valuePart = content.substring(colonIndex + 1).trim();
        
        let highlightedKey = keyPart.replace(REGEX_PATTERNS.JSON_KEY_QUOTES, (match, p1) => HTML_COMPONENTS.JSON_KEY(`"${p1}"`));
        
        let highlightedValue = highlightValue(valuePart);
        
        return leadingSpaces + highlightedKey + HTML_COMPONENTS.JSON_BRACKET(':') + ' ' + highlightedValue;
    }
    
    return leadingSpaces + highlightValue(content);
}

function highlightValue(value) {
    let hasComma = value.endsWith(',');
    if (hasComma) {
        value = value.slice(0, -1).trim();
    }
    
    let result = '';
    
    if (value === 'true' || value === 'false') {
        result = HTML_COMPONENTS.JSON_BOOLEAN(value);
    } else if (value === 'null') {
        result = HTML_COMPONENTS.JSON_NULL(value);
    } else if (REGEX_PATTERNS.NUMERIC_VALUE.test(value)) {
        result = HTML_COMPONENTS.JSON_NUMBER(value);
    } else if (value.startsWith('"') && value.endsWith('"')) {
        result = HTML_COMPONENTS.JSON_STRING(value);
    } else if (value === '[' || value === ']' || value === '{' || value === '}') {
        result = HTML_COMPONENTS.JSON_BRACKET(value);
    } else {
        result = value;
    }
    
    if (hasComma) {
        result += HTML_COMPONENTS.JSON_COMMA();
    }
    
    return result;
}

function copyJsonInput() {
    const element = document.getElementById(ELEMENT_IDS.JSON_INPUT);
    const text = element.value;
    
    if (!text.trim()) {
        window.showError(ERROR_MESSAGES.NO_INPUT_TO_COPY);
        return;
    }
    
    navigator.clipboard.writeText(text).then(() => {
        showJsonNotification(SUCCESS_MESSAGES.INPUT_COPIED);
    }).catch(() => {
        element.select();
        document.execCommand('copy');
        showJsonNotification(SUCCESS_MESSAGES.INPUT_COPIED);
    });
}

function copyJsonOutput() {
    const element = document.getElementById(ELEMENT_IDS.JSON_OUTPUT);
    const text = element.textContent || element.innerText;
    
    if (!text.trim() || text.trim() === UI_TEXT.FORMATTED_JSON_PLACEHOLDER) {
        window.showError(ERROR_MESSAGES.NO_OUTPUT_TO_COPY);
        return;
    }
    
    navigator.clipboard.writeText(text).then(() => {
        showJsonNotification(SUCCESS_MESSAGES.OUTPUT_COPIED);
    }).catch(() => {
        const tempTextarea = document.createElement('textarea');
        tempTextarea.value = text;
        document.body.appendChild(tempTextarea);
        tempTextarea.select();
        document.execCommand('copy');
        document.body.removeChild(tempTextarea);
        showJsonNotification(SUCCESS_MESSAGES.OUTPUT_COPIED);
    });
}

function downloadJSON() {
    const outputElement = document.getElementById(ELEMENT_IDS.JSON_OUTPUT);
    const jsonContent = outputElement.textContent || outputElement.innerText;
    
    if (!jsonContent.trim() || jsonContent.trim() === UI_TEXT.FORMATTED_JSON_PLACEHOLDER) {
        window.showError(ERROR_MESSAGES.NO_JSON_CONTENT_TO_DOWNLOAD);
        return;
    }
    
    const blob = new Blob([jsonContent], { type: DOWNLOAD_CONFIG.BLOB_TYPE });
    
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    
    const now = new Date();
    const timestamp = now.getFullYear() + '-' + 
                    String(now.getMonth() + 1).padStart(2, '0') + '-' + 
                    String(now.getDate()).padStart(2, '0') + '_' + 
                    String(now.getHours()).padStart(2, '0') + '-' + 
                    String(now.getMinutes()).padStart(2, '0') + '-' + 
                    String(now.getSeconds()).padStart(2, '0');
    
    a.download = DOWNLOAD_CONFIG.FILENAME_PREFIX + timestamp + DOWNLOAD_CONFIG.FILENAME_EXTENSION;
    
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    
    window.URL.revokeObjectURL(url);
    
    showJsonNotification(SUCCESS_MESSAGES.JSON_DOWNLOADED);
}

function showJsonNotification(message) {
    const notification = document.createElement('div');
    notification.className = CSS_CLASSES.JSON_NOTIFICATION;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.classList.add(CSS_CLASSES.NOTIFICATION_SHOW);
    }, TIMEOUTS.NOTIFICATION_SHOW);
    
    setTimeout(() => {
        notification.classList.remove(CSS_CLASSES.NOTIFICATION_SHOW);
        notification.classList.add(CSS_CLASSES.NOTIFICATION_HIDE);
        setTimeout(() => notification.remove(), TIMEOUTS.NOTIFICATION_HIDE);
    }, TIMEOUTS.NOTIFICATION_DURATION);
}