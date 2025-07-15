package main

import (
	"Go-Utilities/internal/handlers"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	router := handlers.SetupRoutes()

	port := ":8080"
	url := "http://localhost" + port

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", url)
		if err := http.ListenAndServe(port, router); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Open browser automatically
	openBrowser(url)

	// Keep main thread alive
	select {}
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
