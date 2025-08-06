package main

import (
	"Go-Utilities/internal/consts"
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

	port := consts.DEFAULT_PORT
	url := consts.BASE_URL + port

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		log.Printf(consts.LOG_SERVER_STARTING, url)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(consts.ERR_SERVER_START, err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	openBrowser(url)

	setupGracefulShutdown(server)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		log.Printf(consts.LOG_UNSUPPORTED_PLATFORM, url)
		return
	}

	if err != nil {
		log.Printf(consts.ERR_OPEN_BROWSER, err)
		log.Printf(consts.LOG_OPEN_MANUALLY, url)
	} else {
		log.Printf(consts.LOG_OPENING_BROWSER, url)
	}
}

func setupGracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf(consts.LOG_RECEIVED_SIGNAL, sig)

	handlers.SendShutdownSignal()

	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf(consts.ERR_FORCED_SHUTDOWN, err)
	} else {
		log.Println(consts.LOG_SHUTDOWN_COMPLETE)
	}
}
