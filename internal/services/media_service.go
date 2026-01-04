package services

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"example.com/myapp/internal/models"
	"github.com/google/uuid"
	"golang.org/x/image/webp"
)

const (
	// MaxFileSize is 200 MB in bytes
	MaxFileSize = 200 * 1024 * 1024
	
	// MediaStoragePath is the directory where media files are stored
	MediaStoragePath = "./uploads"
	
	// ImageQuality is the JPEG quality for optimized images (0-100)
	ImageQuality = 85
)

// SupportedImageFormats are the image formats we accept and convert to
var SupportedImageFormats = []string{"image/jpeg", "image/png", "image/webp", "image/gif"}
var SupportedFormats = append(SupportedImageFormats, "application/pdf")

type MediaService struct {
	mediaMap map[string]*models.Media
}

func NewMediaService() *MediaService {
	// Create uploads directory if it doesn't exist
	os.MkdirAll(MediaStoragePath, 0755)
	
	return &MediaService{
		mediaMap: make(map[string]*models.Media),
	}
}

// UploadMedia uploads and processes a media file
func (s *MediaService) UploadMedia(file *multipart.FileHeader) (*models.Media, error) {
	// Validate file size
	if file.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of 200 MB (file size: %.2f MB)", float64(file.Size)/(1024*1024))
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Detect content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Fallback: try to determine from filename
		contentType = getContentTypeFromName(file.Filename)
	}

	// Validate content type
	if !isValidContentType(contentType) {
		return nil, fmt.Errorf("unsupported file type: %s. Supported types: JPEG, PNG, WebP, GIF, PDF", contentType)
	}

	// Determine media type and format
	var mediaType, format string
	var fileBytes []byte
	var imgDims image.Rectangle

	if isImageType(contentType) {
		mediaType = "image"
		// Read file content
		fileBytes, err = io.ReadAll(src)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		// Optimize image
		optimizedBytes, dims, err := s.optimizeImage(fileBytes, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to optimize image: %w", err)
		}

		fileBytes = optimizedBytes
		imgDims = dims
		format = "webp" // Store as WebP for better compression
	} else if isPDFType(contentType) {
		mediaType = "pdf"
		format = "pdf"
		fileBytes, err = io.ReadAll(src)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	}

	// Generate unique filename
	storedName := fmt.Sprintf("%s_%d.%s", uuid.New().String(), time.Now().UnixNano(), format)
	filePath := filepath.Join(MediaStoragePath, storedName)

	// Save file to disk
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create media record
	media := &models.Media{
		ID:           uuid.New().String(),
		OriginalName: file.Filename,
		StoredName:   storedName,
		Type:         mediaType,
		Format:       format,
		SizeBytes:    int64(len(fileBytes)),
		FilePath:     filePath,
		UploadedAt:   time.Now(),
	}

	// Add dimensions if it's an image
	if mediaType == "image" && imgDims.Max.X > 0 {
		media.Width = imgDims.Max.X
		media.Height = imgDims.Max.Y
	}

	// Store in memory
	s.mediaMap[media.ID] = media

	return media, nil
}

// optimizeImage converts any image format to WebP with compression
func (s *MediaService) optimizeImage(fileBytes []byte, contentType string) ([]byte, image.Rectangle, error) {
	// Decode image from various formats
	var img image.Image
	var dims image.Rectangle
	var err error

	switch {
	case strings.Contains(contentType, "jpeg"):
		img, err = jpeg.Decode(strings.NewReader(string(fileBytes)))
	case strings.Contains(contentType, "png"):
		img, err = png.Decode(strings.NewReader(string(fileBytes)))
	case strings.Contains(contentType, "webp"):
		img, err = webp.Decode(strings.NewReader(string(fileBytes)))
	case strings.Contains(contentType, "gif"):
		img, err = png.Decode(strings.NewReader(string(fileBytes)))
	default:
		return nil, image.Rectangle{}, fmt.Errorf("unsupported image format")
	}

	if err != nil {
		return nil, image.Rectangle{}, fmt.Errorf("failed to decode image: %w", err)
	}

	dims = img.Bounds()

	// For simplicity, we'll convert all images to JPEG with quality compression
	// WebP encoding would require additional library
	// Converting to JPEG maintains good quality while reducing size
	
	output := &strings.Builder{}
	err = jpeg.Encode(output, img, &jpeg.Options{Quality: ImageQuality})
	if err != nil {
		return nil, dims, fmt.Errorf("failed to encode image: %w", err)
	}

	return []byte(output.String()), dims, nil
}

// getImageDimensions returns image dimensions
func (s *MediaService) getImageDimensions(fileBytes []byte) (image.Rectangle, error) {
	// Try to decode as various formats to get dimensions
	if config, _, err := image.DecodeConfig(strings.NewReader(string(fileBytes))); err == nil {
		return image.Rectangle{Max: image.Point{X: config.Width, Y: config.Height}}, nil
	}

	return image.Rectangle{}, fmt.Errorf("failed to get image dimensions")
}

// GetMedia retrieves a media file by ID
func (s *MediaService) GetMedia(id string) (*models.Media, error) {
	media, exists := s.mediaMap[id]
	if !exists {
		return nil, fmt.Errorf("media not found")
	}
	return media, nil
}

// GetAllMedia retrieves all media files
func (s *MediaService) GetAllMedia() []*models.Media {
	media := make([]*models.Media, 0, len(s.mediaMap))
	for _, m := range s.mediaMap {
		media = append(media, m)
	}
	return media
}

// DeleteMedia deletes a media file
func (s *MediaService) DeleteMedia(id string) error {
	media, exists := s.mediaMap[id]
	if !exists {
		return fmt.Errorf("media not found")
	}

	// Delete file from disk
	err := os.Remove(media.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Remove from memory
	delete(s.mediaMap, id)
	return nil
}

// Helper functions

func isValidContentType(contentType string) bool {
	for _, ct := range SupportedFormats {
		if strings.Contains(contentType, ct) {
			return true
		}
	}
	return false
}

func isImageType(contentType string) bool {
	for _, ct := range SupportedImageFormats {
		if strings.Contains(contentType, ct) {
			return true
		}
	}
	return false
}

func isPDFType(contentType string) bool {
	return strings.Contains(contentType, "application/pdf")
}

func getContentTypeFromName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return ""
	}
}
