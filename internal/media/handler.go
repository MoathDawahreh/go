package media

import (
	"encoding/json"
	"log/slog"
	"net/http"

	appErr "example.com/myapp/internal/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all media-related routes with appropriate middleware
// Middleware is passed as parameters to avoid circular imports
func (h *Handler) RegisterRoutes(r chi.Router, loggingMw, authMw, validateIDMw func(http.Handler) http.Handler) {
	r.Route("/media", func(r chi.Router) {
		// Middleware for ALL media routes
		r.Use(loggingMw)
		r.Use(authMw)

		r.Post("/upload", h.UploadMedia)
		r.Get("/", h.GetAllMedia)
		r.Route("/{id}", func(r chi.Router) {
			// Middleware for ID-specific operations
			r.Use(validateIDMw)

			r.Get("/", h.GetMedia)
			r.Get("/download", h.DownloadMedia)
			r.Delete("/", h.DeleteMedia)
		})
	})
}

// UploadMedia handles file upload - POST /media/upload
func (h *Handler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 300MB max (buffer for multiple files)
	err := r.ParseMultipartForm(300 * 1024 * 1024)
	if err != nil {
		slog.Error("Failed to parse form", "error", err)
		respondMediaError(w, appErr.BadRequest("failed to parse form"), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to get file from request", "error", err)
		respondMediaError(w, appErr.BadRequest("failed to get file from request"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload and process media
	media, err := h.service.UploadMedia(r.Context(), fileHeader)

	response := &MediaUploadResponse{
		Success: err == nil,
		Media:   media,
	}

	if err != nil {
		response.Error = err.Error()
		response.Success = false
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(getMediaStatusCode(err))
	} else {
		response.Message = "File uploaded and processed successfully"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		slog.Info("Media uploaded", "id", media.ID)
	}

	json.NewEncoder(w).Encode(response)
}

// GetAllMedia retrieves all media files - GET /media
func (h *Handler) GetAllMedia(w http.ResponseWriter, r *http.Request) {
	mediaList, err := h.service.GetAllMedia(r.Context())
	if err != nil {
		slog.Error("Failed to get all media", "error", err)
		response := &MediaListResponse{
			Error: "Failed to retrieve media: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &MediaListResponse{
		Total: len(mediaList),
		Media: mediaList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetMedia retrieves a specific media file - GET /media/{id}
func (h *Handler) GetMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	media, err := h.service.GetMedia(r.Context(), id)
	if err != nil {
		slog.Error("Failed to get media", "id", id, "error", err)
		response := &MediaListResponse{
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
func (h *Handler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.DeleteMedia(r.Context(), id)
	if err != nil {
		slog.Error("Failed to delete media", "id", id, "error", err)
		response := &MediaListResponse{
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

	slog.Info("Media deleted", "id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DownloadMedia serves a media file - GET /media/{id}/download
func (h *Handler) DownloadMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	media, err := h.service.GetMedia(r.Context(), id)
	if err != nil {
		slog.Error("Failed to get media for download", "id", id, "error", err)
		http.Error(w, "Media not found", http.StatusNotFound)
		return
	}

	// Serve the file
	w.Header().Set("Content-Disposition", "attachment; filename="+media.OriginalName)
	w.Header().Set("Content-Type", getMediaContentType(media.Format))
	http.ServeFile(w, r, media.FilePath)
}

// Helper functions

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

func getMediaStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	ae := appErr.GetAppError(err)
	if ae == nil {
		return http.StatusInternalServerError
	}

	switch ae.Code {
	case appErr.ErrCodeNotFound:
		return http.StatusNotFound
	case appErr.ErrCodeBadRequest, appErr.ErrCodeInvalidID, appErr.ErrCodeUnsupported:
		return http.StatusBadRequest
	case appErr.ErrCodeFileTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

func respondMediaError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := &MediaListResponse{
		Error: err.Error(),
	}
	json.NewEncoder(w).Encode(response)
}
