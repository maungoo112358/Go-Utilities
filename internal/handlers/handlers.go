package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
	"Go-Utilities/internal/consts"
	"Go-Utilities/internal/downloader"
	"Go-Utilities/internal/models"
	
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var downloadManager *downloader.Manager
var shutdownSignal = make(chan bool, 10) // Buffered channel for shutdown signals

func init() {
	downloadManager = downloader.NewManager()
	
	if err := downloadManager.TestYtDlp(); err != nil {
		log.Printf(consts.LOG_YT_DLP_TEST_FAILED, err)
	}
}

// SendShutdownSignal sends shutdown signal to all connected WebSocket clients
func SendShutdownSignal() {
	log.Println(consts.LOG_SENDING_SHUTDOWN_SIGNAL)
	select {
	case shutdownSignal <- true:
		log.Println(consts.LOG_SHUTDOWN_SIGNAL_SENT)
	default:
		log.Println(consts.LOG_SHUTDOWN_SIGNAL_FULL)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf(consts.LOG_TEMPLATE_ERROR, err)
		http.Error(w, consts.ERR_TEMPLATE+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	log.Printf("Starting download for URL: %s, Quality: %s", req.URL, req.Quality)
	downloadID := downloadManager.StartDownload(req.URL, req.Quality)
	log.Printf("Download started with ID: %s", downloadID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success: true,
		Message: "Download started",
		FileName: downloadID,
	})
}

func VideoInfoHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	videoInfo, err := downloadManager.GetVideoInfo(req.URL)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videoInfo)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	
	log.Printf("WebSocket connection established")
	updates := downloadManager.SubscribeToUpdates()
	
	// Listen for both download updates and shutdown signals
	for {
		select {
		case update := <-updates:
			log.Printf("Sending WebSocket update: %+v", update)
			if err := conn.WriteJSON(update); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
			
		case <-shutdownSignal:
			log.Printf("Sending shutdown signal to WebSocket client")
			shutdownMsg := map[string]interface{}{
				"type": "shutdown",
				"message": "Application is shutting down",
			}
			if err := conn.WriteJSON(shutdownMsg); err != nil {
				log.Printf("Failed to send shutdown signal: %v", err)
			}
			// Give client time to process shutdown signal
			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}

func Mp3ConvertHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	log.Printf("Starting MP3 conversion for URL: %s", req.URL)
	downloadID := downloadManager.StartMp3Convert(req.URL)
	log.Printf("MP3 conversion started with ID: %s", downloadID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success: true,
		Message: "MP3 conversion started",
		FileName: downloadID,
	})
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	// Return JavaScript that closes the current tab
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Shutting Down</title>
    <style>
        body {
            font-family: 'JetBrains Mono', monospace;
            background-color: #121212;
            color: #FFFFFF;
            text-align: center;
            padding: 50px;
        }
    </style>
</head>
<body>
    <h1>Application Shutting Down</h1>
    <p>This tab will close automatically.</p>
    <script>
        // Try multiple methods to close the tab
        setTimeout(() => {
            // Method 1: window.close()
            window.close();
            
            // Method 2: If window.close() doesn't work, try to navigate away
            setTimeout(() => {
                window.location.href = 'about:blank';
            }, 500);
        }, 1000);
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success: false,
		Message: message,
	})
}