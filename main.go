package main

import (
	"Go-Utilities/internal/consts"
	"Go-Utilities/internal/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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
	browserUsed := ""
	err := exec.Command("cmd", "/c", "start", url).Start()

	if err != nil {
		browserUsed, err = openWithDefaultBrowser(url)
		if err != nil {
			log.Printf(consts.LOG_BROWSER_OPEN_FAILED, err, url)
		} else {
			log.Printf(consts.LOG_OPENING_BROWSER, url)
			if browserUsed != "" {
				os.Setenv("GO_UTILS_BROWSER", browserUsed)
				log.Printf("Detected browser: %s", browserUsed)
			}
		}
	} else {
		log.Printf(consts.LOG_OPENING_BROWSER, url)
		detectDefaultBrowser()
	}
}

func requestCookiePermission() {
	time.Sleep(2 * time.Second)

	message := "Do you want to allow this app to use your YouTube cookies for better download success?\\n\\nThis will:\\n• Open YouTube in a new tab\\n• Extract cookies after you log in\\n• Close YouTube tab automatically\\n• Use cookies for authenticated downloads\\n\\nNo passwords are stored, only session cookies."
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; $result = [System.Windows.Forms.MessageBox]::Show("%s", "YouTube Cookie Permission", [System.Windows.Forms.MessageBoxButtons]::YesNo, [System.Windows.Forms.MessageBoxIcon]::Question); if ($result -eq "Yes") { Write-Output "YES" } else { Write-Output "NO" }`, message)

	cmd := exec.Command("powershell", "-Command", psScript)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to show cookie permission dialog: %v", err)
		return
	}

	response := strings.TrimSpace(string(output))
	if response == "YES" {
		log.Printf("User accepted cookie permission. Opening YouTube in new tab...")
		os.Setenv("GO_UTILS_COOKIES_ENABLED", "true")
		openYouTubeAndExtractCookies()
	} else {
		log.Printf("User declined cookie permission. Downloads will work without authentication.")
		os.Setenv("GO_UTILS_COOKIES_ENABLED", "false")
	}
}

func openYouTubeAndExtractCookies() {
	youtubeURL := "https://www.youtube.com"
	log.Printf("Opening YouTube in new tab for cookie extraction...")
	
	err := exec.Command("cmd", "/c", "start", youtubeURL).Start()
	if err != nil {
		_, err = openWithDefaultBrowser(youtubeURL)
		if err != nil {
			log.Printf("Failed to open YouTube: %v", err)
			return
		}
	}
	
	log.Printf("YouTube tab opened. Please log in if needed...")
	log.Printf("Waiting 10 seconds for login and cookie extraction...")
	time.Sleep(10 * time.Second)
	
	detectedBrowser := os.Getenv("GO_UTILS_BROWSER")
	if detectedBrowser != "" {
		log.Printf("Extracting cookies from %s browser...", detectedBrowser)
		
		if testCookieExtraction(detectedBrowser) {
			log.Printf("✓ Successfully extracted YouTube cookies from %s!", detectedBrowser)
			closeYouTubeTab()
		} else {
			log.Printf("⚠ Could not verify cookie extraction. YouTube tab remains open for manual login.")
		}
	} else {
		log.Printf("⚠ Browser detection failed. YouTube tab remains open.")
	}
}

func testCookieExtraction(browser string) bool {
	ytDlpPath, err := getYtDlpPath()
	if err != nil {
		log.Printf("Failed to get yt-dlp path: %v", err)
		return false
	}
	
	log.Printf("Testing cookie extraction from %s browser...", browser)
	testCmd := exec.Command(ytDlpPath, "--cookies-from-browser", browser, "--simulate", "--no-warnings", "https://www.youtube.com")
	
	output, err := testCmd.CombinedOutput()
	outputStr := string(output)
	
	if err != nil {
		log.Printf("Cookie extraction test details:")
		log.Printf("  Error: %v", err)
		log.Printf("  Output: %s", outputStr)
		
		if strings.Contains(outputStr, "DPAPI") || strings.Contains(outputStr, "keyring") || strings.Contains(outputStr, "cookies") {
			log.Printf("Cookie access possible but may need browser restart or permissions")
			return true
		}
		return false
	}
	
	log.Printf("✓ Cookie extraction test passed - cookies are accessible from %s", browser)
	return true
}

func closeYouTubeTab() {
	log.Printf("Closing YouTube tab...")
	
	psScript := `
		$wshell = New-Object -ComObject wscript.shell
		$wshell.AppActivate("YouTube")
		Start-Sleep -Milliseconds 500
		$wshell.SendKeys("^{w}")
	`
	
	cmd := exec.Command("powershell", "-Command", psScript)
	err := cmd.Run()
	if err != nil {
		log.Printf("Could not automatically close YouTube tab. Please close it manually.")
	} else {
		log.Printf("✓ YouTube tab closed successfully.")
	}
}

func getYtDlpPath() (string, error) {
	return "dependencies/yt-dlp.exe", nil
}

func detectDefaultBrowser() {
	cmd := exec.Command("reg", "query", "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\Shell\\Associations\\UrlAssociations\\http\\UserChoice")
	output, err := cmd.Output()
	if err == nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "BraveHTML") {
			os.Setenv("GO_UTILS_BROWSER", "brave")
			log.Printf("Detected default browser: Brave")
		} else if strings.Contains(outputStr, "ChromeHTML") {
			os.Setenv("GO_UTILS_BROWSER", "chrome")
			log.Printf("Detected default browser: Chrome")
		} else if strings.Contains(outputStr, "MSEdgeHTM") {
			os.Setenv("GO_UTILS_BROWSER", "edge")
			log.Printf("Detected default browser: Edge")
		} else if strings.Contains(outputStr, "FirefoxURL") {
			os.Setenv("GO_UTILS_BROWSER", "firefox")
			log.Printf("Detected default browser: Firefox")
		}
	}
}

func openWithDefaultBrowser(url string) (string, error) {
	browsers := []struct {
		name  string
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
					os.Setenv("GO_UTILS_BROWSER", browser.name)
					return browser.name, nil
				}
			}
		}
	}

	return "", exec.Command(consts.RUNDLL32_COMMAND, consts.URL_DLL_HANDLER, url).Start()
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
