package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
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

func init() {
	downloadManager = downloader.NewManager()
	
	// Test yt-dlp on startup
	if err := downloadManager.TestYtDlp(); err != nil {
		log.Printf("WARNING: yt-dlp test failed: %v", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
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
	for update := range updates {
		log.Printf("Sending WebSocket update: %+v", update)
		if err := conn.WriteJSON(update); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
	log.Printf("WebSocket connection closed")
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success: false,
		Message: message,
	})
}