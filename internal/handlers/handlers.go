package handlers

import (
	"Go-Utilities/internal/consts"
	"Go-Utilities/internal/downloader"
	"Go-Utilities/internal/models"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var downloadManager *downloader.Manager
var shutdownSignal = make(chan bool, consts.SHUTDOWN_SIGNAL_BUFFER) // Buffered channel for shutdown signals

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
	tmpl, err := template.ParseFiles(consts.TEMPLATE_PATH)
	if err != nil {
		log.Printf(consts.LOG_TEMPLATE_ERROR, err)
		http.Error(w, consts.ERR_TEMPLATE+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf(consts.LOG_TEMPLATE_EXECUTION_ERROR, err)
		http.Error(w, consts.ERR_TEMPLATE_EXECUTION, http.StatusInternalServerError)
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf(consts.LOG_INVALID_REQUEST_BODY, err)
		sendJSONError(w, consts.ERR_INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	log.Printf(consts.LOG_STARTING_DOWNLOAD, req.URL, req.Quality)
	downloadID := downloadManager.StartDownload(req.URL, req.Quality)
	log.Printf(consts.LOG_DOWNLOAD_STARTED, downloadID)

	w.Header().Set(consts.HEADER_CONTENT_TYPE, consts.CONTENT_TYPE_JSON)
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success:  true,
		Message:  consts.MSG_DOWNLOAD_STARTED,
		FileName: downloadID,
	})
}

func VideoInfoHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, consts.ERR_INVALID_REQUEST_INFO, http.StatusBadRequest)
		return
	}

	videoInfo, err := downloadManager.GetVideoInfo(req.URL)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(consts.HEADER_CONTENT_TYPE, consts.CONTENT_TYPE_JSON)
	json.NewEncoder(w).Encode(videoInfo)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf(consts.LOG_WS_UPGRADE_ERROR, err)
		return
	}
	defer conn.Close()

	log.Printf(consts.LOG_WS_CONNECTION_ESTABLISHED)
	updates := downloadManager.SubscribeToUpdates()

	// Listen for both download updates and shutdown signals
	for {
		select {
		case update := <-updates:
			log.Printf(consts.LOG_SENDING_WS_UPDATE, update)
			if err := conn.WriteJSON(update); err != nil {
				log.Printf(consts.LOG_WS_WRITE_ERROR, err)
				return
			}

		case <-shutdownSignal:
			log.Printf(consts.LOG_SENDING_SHUTDOWN_TO_WS)
			shutdownMsg := map[string]interface{}{
				"type":    consts.WS_MESSAGE_TYPE_SHUTDOWN,
				"message": consts.MSG_SHUTDOWN_SIGNAL,
			}
			if err := conn.WriteJSON(shutdownMsg); err != nil {
				log.Printf(consts.ERR_SEND_SHUTDOWN_SIGNAL, err)
			}
			// Give client time to process shutdown signal
			time.Sleep(consts.SHUTDOWN_DELAY_MS * time.Millisecond)
			return
		}
	}
}

func Mp3ConvertHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf(consts.LOG_INVALID_REQUEST_BODY_MP3, err)
		sendJSONError(w, consts.ERR_INVALID_REQUEST_MP3, http.StatusBadRequest)
		return
	}

	log.Printf(consts.LOG_STARTING_MP3_CONVERSION, req.URL)
	downloadID := downloadManager.StartMp3Convert(req.URL)
	log.Printf(consts.LOG_MP3_CONVERSION_STARTED, downloadID)

	w.Header().Set(consts.HEADER_CONTENT_TYPE, consts.CONTENT_TYPE_JSON)
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success:  true,
		Message:  consts.MSG_MP3_CONVERSION_STARTED,
		FileName: downloadID,
	})
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(consts.SHUTDOWN_TEMPLATE_PATH)
	if err != nil {
		log.Printf(consts.LOG_TEMPLATE_ERROR, err)
		http.Error(w, consts.ERR_TEMPLATE+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf(consts.LOG_TEMPLATE_EXECUTION_ERROR, err)
		http.Error(w, consts.ERR_TEMPLATE_EXECUTION, http.StatusInternalServerError)
	}
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set(consts.HEADER_CONTENT_TYPE, consts.CONTENT_TYPE_JSON)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.DownloadResponse{
		Success: false,
		Message: message,
	})
}
