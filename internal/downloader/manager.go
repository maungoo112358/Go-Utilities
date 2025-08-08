package downloader

import (
	"Go-Utilities/internal/consts"
	"Go-Utilities/internal/models"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	downloads   map[string]*Download
	subscribers []chan models.ProgressUpdate
	mu          sync.RWMutex
}

type Download struct {
	ID       string
	URL      string
	Title    string
	Status   string
	Progress float64
	Speed    string
	ETA      string
}

func NewManager() *Manager {
	return &Manager{
		downloads:   make(map[string]*Download),
		subscribers: []chan models.ProgressUpdate{},
	}
}

func (m *Manager) TestYtDlp() error {
	return TestYtDlp()
}

func (m *Manager) StartDownload(url, quality string) string {
	downloadID := fmt.Sprintf(consts.DOWNLOAD_ID_FORMAT, time.Now().Unix())
	go m.download(downloadID, url, quality)
	return downloadID
}

func (m *Manager) StartMp3Convert(url string) string {
	downloadID := fmt.Sprintf(consts.MP3_ID_FORMAT, time.Now().Unix())
	go m.convertToMp3(downloadID, url)
	return downloadID
}

func (m *Manager) download(id, url, quality string) {
	m.mu.Lock()
	m.downloads[id] = &Download{
		ID:     id,
		URL:    url,
		Status: consts.STATUS_STARTING,
	}
	m.mu.Unlock()

	m.updateStatus(id, consts.STATUS_DOWNLOADING, 0, "", "", consts.MSG_STARTING_DOWNLOAD)

	result, err := ExecuteDownload(url, quality, func(progress float64, speed, eta, message string) {
		m.updateStatus(id, consts.STATUS_DOWNLOADING, progress, speed, eta, message)
	})

	if err != nil {
		m.updateStatus(id, consts.STATUS_ERROR, 0, "", "", err.Error())
		return
	}

	if result.Title != "" {
		m.mu.Lock()
		if download, ok := m.downloads[id]; ok {
			download.Title = result.Title
		}
		m.mu.Unlock()
	}

	newFileName := m.addResolutionToFilename(result.FilePath, quality)
	if newFileName != result.FilePath {
		if err := os.Rename(result.FilePath, newFileName); err != nil {
			log.Printf(consts.ERR_RENAME_FILE, err)
		} else {
			result.FilePath = newFileName
		}
	}

	finalPath, err := m.openFilePicker(result.FilePath)
	if err != nil {
		m.updateStatus(id, consts.STATUS_ERROR, 0, "", "", fmt.Sprintf(consts.ERR_SAVE_FILE, err))
		return
	}

	m.updateStatus(id, consts.STATUS_COMPLETED, 100, "", "", fmt.Sprintf(consts.MSG_SAVED_AS, filepath.Base(finalPath)))

	log.Printf(consts.LOG_OPENING_FILE_EXPLORER)
	m.openFileExplorer(filepath.Dir(finalPath))
}

func (m *Manager) convertToMp3(id, url string) {
	m.mu.Lock()
	m.downloads[id] = &Download{
		ID:     id,
		URL:    url,
		Status: consts.STATUS_STARTING,
	}
	m.mu.Unlock()

	m.updateStatus(id, consts.STATUS_CONVERTING, 0, "", "", consts.MSG_STARTING_MP3_CONVERSION)

	result, err := ExecuteMp3Conversion(url, func(progress float64, speed, eta, message string) {
		m.updateStatus(id, consts.STATUS_CONVERTING, progress, speed, eta, message)
	})

	if err != nil {
		m.updateStatus(id, consts.STATUS_ERROR, 0, "", "", fmt.Sprintf(consts.MP3_CONVERSION_FAILED, err))
		return
	}

	if result == nil {
		m.updateStatus(id, consts.STATUS_ERROR, 0, "", "", "MP3 conversion failed: no result returned")
		return
	}

	if result.Title != "" {
		m.mu.Lock()
		if download, ok := m.downloads[id]; ok {
			download.Title = result.Title
		}
		m.mu.Unlock()
	}

	finalPath, err := m.openFilePicker(result.FilePath)
	if err != nil {
		m.updateStatus(id, consts.STATUS_ERROR, 0, "", "", fmt.Sprintf(consts.ERR_SAVE_MP3_FILE, err))
		return
	}

	m.updateStatus(id, consts.STATUS_COMPLETED, 100, "", "", fmt.Sprintf(consts.MSG_MP3_SAVED_AS, filepath.Base(finalPath)))
}

func (m *Manager) GetVideoInfo(url string) (*models.VideoInfo, error) {
	return GetVideoInfo(url)
}

func (m *Manager) ParseYouTubeURL(inputURL string) (string, error) {
	return ParseYouTubeURL(inputURL)
}

func (m *Manager) updateStatus(id, status string, progress float64, speed, eta, message string) {
	m.mu.Lock()
	if download, ok := m.downloads[id]; ok {
		download.Status = status
		download.Progress = progress
		download.Speed = speed
		download.ETA = eta
	}
	m.mu.Unlock()

	update := models.ProgressUpdate{
		ID:       id,
		Progress: progress,
		Speed:    speed,
		ETA:      eta,
		Status:   status,
		Message:  message,
	}

	log.Printf(consts.LOG_BROADCASTING_UPDATE, update)
	m.broadcast(update)
}

func (m *Manager) broadcast(update models.ProgressUpdate) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	log.Printf(consts.LOG_BROADCASTING_SUBSCRIBERS, len(m.subscribers))
	for i, ch := range m.subscribers {
		select {
		case ch <- update:
			log.Printf(consts.LOG_SENT_UPDATE, i)
		default:
			log.Printf(consts.LOG_SUBSCRIBER_CHANNEL_FULL, i)
		}
	}
}

func (m *Manager) SubscribeToUpdates() <-chan models.ProgressUpdate {
	ch := make(chan models.ProgressUpdate, 100)
	m.mu.Lock()
	m.subscribers = append(m.subscribers, ch)
	m.mu.Unlock()
	return ch
}

func (m *Manager) addResolutionToFilename(filePath, quality string) string {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	if quality != "" && quality != consts.BEST_QUALITY && !strings.Contains(nameWithoutExt, quality) {
		newFilename := fmt.Sprintf("%s [%s]%s", nameWithoutExt, quality, ext)
		return filepath.Join(dir, newFilename)
	}

	return filePath
}

func (m *Manager) copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func (m *Manager) openFilePicker(sourceFile string) (string, error) {
	filename := filepath.Base(sourceFile)

	psScript := fmt.Sprintf(consts.POWERSHELL_FILE_PICKER_SCRIPT,
		consts.FILE_PICKER_FILTER,
		filename,
		consts.SAVE_DIALOG_TITLE,
		filepath.Ext(filename),
		consts.CANCELLED_MESSAGE)

	cmd := exec.Command(consts.POWERSHELL_COMMAND, consts.COMMAND_FLAG, psScript)
	output, err := cmd.Output()
	if err != nil {
		log.Printf(consts.POWERSHELL_FAILED, err)
		return "", fmt.Errorf(consts.ERR_SAVE_FILE_PICKER, err)
	}

	selectedPath := strings.TrimSpace(string(output))
	if selectedPath == consts.CANCELLED_MESSAGE || selectedPath == "" {
		os.Remove(sourceFile)
		return "", fmt.Errorf(consts.ERR_SAVE_CANCELLED)
	}

	if err := m.copyFile(sourceFile, selectedPath); err != nil {
		return "", fmt.Errorf(consts.ERR_SAVE_FILE_PICKER, err)
	}

	os.Remove(sourceFile)
	return selectedPath, nil
}

func (m *Manager) openFileExplorer(path string) {
	exec.Command(consts.EXPLORER_COMMAND, path).Start()
}
