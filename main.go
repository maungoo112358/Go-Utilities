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
	err := exec.Command("cmd", "/c", "start", url).Start()
	
	if err != nil {
		err = openWithDefaultBrowser(url)
		if err != nil {
			log.Printf(consts.LOG_BROWSER_OPEN_FAILED, err, url)
		} else {
			log.Printf(consts.LOG_OPENING_BROWSER, url)
		}
	} else {
		log.Printf(consts.LOG_OPENING_BROWSER, url)
	}
}

func openWithDefaultBrowser(url string) error {
	browsers := []struct {
		name string
		paths []string
	}{
		{
			name: consts.BROWSER_NAME_CHROME,
			paths: []string{
				consts.CHROME_PATH_1,
				consts.CHROME_PATH_2,
				os.Getenv("LOCALAPPDATA") + consts.CHROME_PATH_USER,
			},
		},
		{
			name: consts.BROWSER_NAME_FIREFOX, 
			paths: []string{
				consts.FIREFOX_PATH_1,
				consts.FIREFOX_PATH_2,
			},
		},
		{
			name: consts.BROWSER_NAME_EDGE,
			paths: []string{
				consts.EDGE_PATH_1,
				consts.EDGE_PATH_2,
			},
		},
		{
			name: consts.BROWSER_NAME_BRAVE,
			paths: []string{
				consts.BRAVE_PATH_1,
				consts.BRAVE_PATH_2,
			},
		},
	}

	for _, browser := range browsers {
		for _, path := range browser.paths {
			if _, err := os.Stat(path); err == nil {
				if err := exec.Command(path, url).Start(); err == nil {
					return nil
				}
			}
		}
	}

	return exec.Command(consts.RUNDLL32_COMMAND, consts.URL_DLL_HANDLER, url).Start()
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
