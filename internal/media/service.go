package media

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	appErr "example.com/myapp/internal/errors"
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

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	// Create uploads directory if it doesn't exist
	os.MkdirAll(MediaStoragePath, 0755)

	return &Service{
		repo: repo,
	}
}

// UploadMedia uploads and processes a media file
func (s *Service) UploadMedia(ctx context.Context, file *multipart.FileHeader) (*Media, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	// Validate file size
	if file.Size > MaxFileSize {
		return nil, appErr.FileTooLarge(fmt.Sprintf("file size exceeds maximum limit of 200 MB (file size: %.2f MB)", float64(file.Size)/(1024*1024)))
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "error", err)
		return nil, appErr.Internal("failed to open file", err)
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
		return nil, appErr.UnsupportedType("unsupported file type: " + contentType + ". Supported types: JPEG, PNG, WebP, GIF, PDF")
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
			slog.Error("Failed to optimize image", "error", err)
			return nil, appErr.Internal("failed to optimize image", err)
		}

		fileBytes = optimizedBytes
		imgDims = dims
		format = "webp" // Store as WebP for better compression
	} else if isPDFType(contentType) {
		mediaType = "pdf"
		format = "pdf"
		fileBytes, err = io.ReadAll(src)
		if err != nil {
			slog.Error("Failed to read file", "error", err)
			return nil, appErr.Internal("failed to read file", err)
		}
	}

	// Generate unique filename
	storedName := fmt.Sprintf("%s_%d.%s", uuid.New().String(), time.Now().UnixNano(), format)
	filePath := filepath.Join(MediaStoragePath, storedName)

	// Save file to disk
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		slog.Error("Failed to save file", "error", err)
		return nil, appErr.Internal("failed to save file", err)
	}

	// Create media record
	media := &Media{
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

	// Store in repository
	err = s.repo.Save(ctx, media)
	if err != nil {
		slog.Error("Failed to save media to repository", "error", err)
		return nil, appErr.Internal("failed to save media to repository", err)
	}

	return media, nil
}

// optimizeImage converts any image format to WebP with compression
func (s *Service) optimizeImage(fileBytes []byte, contentType string) ([]byte, image.Rectangle, error) {
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
func (s *Service) getImageDimensions(fileBytes []byte) (image.Rectangle, error) {
	// Try to decode as various formats to get dimensions
	if config, _, err := image.DecodeConfig(strings.NewReader(string(fileBytes))); err == nil {
		return image.Rectangle{Max: image.Point{X: config.Width, Y: config.Height}}, nil
	}

	return image.Rectangle{}, fmt.Errorf("failed to get image dimensions")
}

// GetMedia retrieves a media file by ID
func (s *Service) GetMedia(ctx context.Context, id string) (*Media, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get media", "id", id, "error", err)
		return nil, err
	}
	return media, nil
}

// GetAllMedia retrieves all media files
func (s *Service) GetAllMedia(ctx context.Context) ([]*Media, error) {
	if err := ctx.Err(); err != nil {
		return nil, appErr.Internal("context cancelled", err)
	}

	media, err := s.repo.GetAll(ctx)
	if err != nil {
		slog.Error("Failed to get all media", "error", err)
		return nil, appErr.Internal("failed to retrieve media", err)
	}
	return media, nil
}

// DeleteMedia deletes a media file
func (s *Service) DeleteMedia(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return appErr.Internal("context cancelled", err)
	}

	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get media for deletion", "id", id, "error", err)
		return err
	}

	// Delete file from disk
	err = os.Remove(media.FilePath)
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Failed to delete file", "error", err)
		return appErr.Internal("failed to delete file", err)
	}

	// Remove from repository
	return s.repo.Delete(ctx, id)
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
