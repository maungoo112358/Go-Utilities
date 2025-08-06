package handlers

import (
	"Go-Utilities/internal/consts"
	"net/http"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	// Static files
	fs := http.FileServer(http.Dir(consts.STATIC_DIR_PATH))
	r.PathPrefix(consts.STATIC_ROUTE_PREFIX).Handler(http.StripPrefix(consts.STATIC_ROUTE_PREFIX, fs))
	
	// Main page
	r.HandleFunc(consts.HOME_ROUTE, HomeHandler).Methods(consts.HTTP_GET)
	
	// Shutdown page (for graceful browser closure)
	r.HandleFunc(consts.SHUTDOWN_ROUTE, ShutdownHandler).Methods(consts.HTTP_GET)
	
	// API routes
	api := r.PathPrefix(consts.API_ROUTE_PREFIX).Subrouter()
	api.HandleFunc(consts.DOWNLOAD_ROUTE, DownloadHandler).Methods(consts.HTTP_POST)
	api.HandleFunc(consts.MP3_CONVERT_ROUTE, Mp3ConvertHandler).Methods(consts.HTTP_POST)
	api.HandleFunc(consts.VIDEO_INFO_ROUTE, VideoInfoHandler).Methods(consts.HTTP_POST)
	api.HandleFunc(consts.WEBSOCKET_ROUTE, WebSocketHandler)
	
	return r
}