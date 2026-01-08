package media

import "time"

// Media represents a stored media file
type Media struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	Type         string    `json:"type"` // image, pdf
	Format       string    `json:"format"` // jpg, png, webp, pdf
	SizeBytes    int64     `json:"size_bytes"`
	FilePath     string    `json:"file_path"`
	UploadedAt   time.Time `json:"uploaded_at"`
	Width        int       `json:"width,omitempty"` // For images
	Height       int       `json:"height,omitempty"` // For images
}

// MediaUploadResponse is the response after uploading media
type MediaUploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Media   *Media `json:"media,omitempty"`
	Error   string `json:"error,omitempty"`
}

// MediaListResponse is the response when listing media
type MediaListResponse struct {
	Total  int      `json:"total"`
	Media  []*Media `json:"media"`
	Error  string   `json:"error,omitempty"`
}
