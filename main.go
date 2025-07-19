package main

import (
	"Go-Utilities/internal/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	router := handlers.SetupRoutes()

	port := ":8484"
	url := "http://localhost" + port

	// Create HTTP server
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", url)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Open browser automatically
	openBrowser(url)

	// Set up graceful shutdown
	setupGracefulShutdown(server, url)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		// Windows: start command opens in existing browser or launches new one
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		// macOS
		err = exec.Command("open", url).Start()
	case "linux":
		// Linux
		err = exec.Command("xdg-open", url).Start()
	default:
		log.Printf("Unsupported platform, please open %s manually", url)
		return
	}

	if err != nil {
		log.Printf("Failed to open browser automatically: %v", err)
		log.Printf("Please open %s manually", url)
	} else {
		log.Printf("Opening %s in your default browser...", url)
	}
}

func setupGracefulShutdown(server *http.Server, url string) {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)

	// Register the channel to receive specific signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Send shutdown signal to WebSocket clients first
	handlers.SendShutdownSignal()

	// Give clients time to receive shutdown signal and close
	time.Sleep(1 * time.Second)

	// Create a context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown complete")
	}
}
