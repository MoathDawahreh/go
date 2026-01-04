package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/myapp/internal/models"
	"example.com/myapp/internal/services"
)

type MediaHandler struct {
	service *services.MediaService
}

func NewMediaHandler(service *services.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

// UploadMedia handles file upload - POST /media/upload
func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 300MB max (buffer for multiple files)
	err := r.ParseMultipartForm(300 * 1024 * 1024)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload and process media
	media, err := h.service.UploadMedia(fileHeader)
	
	response := &models.MediaUploadResponse{
		Success: err == nil,
		Media:   media,
	}

	if err != nil {
		response.Error = err.Error()
		response.Success = false
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		response.Message = "File uploaded and processed successfully"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(response)
}

// GetAllMedia retrieves all media files - GET /media
func (h *MediaHandler) GetAllMedia(w http.ResponseWriter, r *http.Request) {
	mediaList := h.service.GetAllMedia()

	response := &models.MediaListResponse{
		Total: len(mediaList),
		Media: mediaList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetMedia retrieves a specific media file - GET /media/{id}
func (h *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	media, err := h.service.GetMedia(id)
	if err != nil {
		response := &models.MediaListResponse{
			Error: "Media not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(media)
}

// DeleteMedia deletes a media file - DELETE /media/{id}
func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.DeleteMedia(id)
	if err != nil {
		response := &models.MediaListResponse{
			Error: "Failed to delete media: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Media deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DownloadMedia serves a media file - GET /media/{id}/download
func (h *MediaHandler) DownloadMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	media, err := h.service.GetMedia(id)
	if err != nil {
		http.Error(w, "Media not found", http.StatusNotFound)
		return
	}

	// Serve the file
	w.Header().Set("Content-Disposition", "attachment; filename="+media.OriginalName)
	w.Header().Set("Content-Type", getMediaContentType(media.Format))
	http.ServeFile(w, r, media.FilePath)
}

// Helper function to get content type based on format
func getMediaContentType(format string) string {
	switch format {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	case "gif":
		return "image/gif"
	case "pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
