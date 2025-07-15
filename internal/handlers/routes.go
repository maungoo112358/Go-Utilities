package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	// Static files
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
	// Main page
	r.HandleFunc("/", HomeHandler).Methods("GET")
	
	// Shutdown page (for graceful browser closure)
	r.HandleFunc("/shutdown", ShutdownHandler).Methods("GET")
	
	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/download", DownloadHandler).Methods("POST")
	api.HandleFunc("/mp3-convert", Mp3ConvertHandler).Methods("POST")
	api.HandleFunc("/video-info", VideoInfoHandler).Methods("POST")
	api.HandleFunc("/ws", WebSocketHandler)
	
	return r
}