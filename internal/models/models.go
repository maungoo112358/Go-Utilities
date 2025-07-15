package models

type DownloadRequest struct {
	URL     string `json:"url"`
	Quality string `json:"quality"`
}

type DownloadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FileName string `json:"filename,omitempty"`
	FilePath string `json:"filepath,omitempty"`
}

type ProgressUpdate struct {
	ID         string  `json:"id"`
	Progress   float64 `json:"progress"`
	Speed      string  `json:"speed"`
	ETA        string  `json:"eta"`
	Status     string  `json:"status"`
	Message    string  `json:"message,omitempty"`
}

type VideoFormat struct {
	FormatID   string `json:"format_id"`
	Resolution string `json:"resolution"`
	Extension  string `json:"ext"`
	FileSize   string `json:"filesize"`
	Quality    string `json:"quality"`
}

type VideoInfo struct {
	Title       string        `json:"title"`
	Duration    string        `json:"duration"`
	Thumbnail   string        `json:"thumbnail"`
	Formats     []VideoFormat `json:"formats"`
	ParsedURL   string        `json:"parsed_url"`
}