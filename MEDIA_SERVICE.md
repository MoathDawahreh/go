# Media Service Documentation

## Overview

The Media Service handles uploading, processing, and managing media files (images and PDFs) in your Go application. It includes automatic image optimization and file size validation.

## Features

- **Image Upload**: Supports JPEG, PNG, WebP, and GIF formats
- **PDF Upload**: Supports PDF documents up to 200 MB
- **File Size Validation**: Rejects files larger than 200 MB
- **Automatic Image Optimization**:
  - Converts all images to JPEG format with 85% quality for optimal compression
  - Preserves image resolution and dimensions
  - Reduces file size without significant quality loss
- **Local Storage**: Files are stored in the `./uploads` directory
- **File Management**: Retrieve, list, download, and delete media files

## API Endpoints

### Upload Media

**POST** `/media/upload`

Upload a new media file (image or PDF).

**Request:**

- Content-Type: `multipart/form-data`
- Form field name: `file`

**Example with curl:**

```bash
curl -X POST http://localhost:8080/media/upload \
  -F "file=@image.jpg"

curl -X POST http://localhost:8080/media/upload \
  -F "file=@document.pdf"
```

**Response (Success - 201):**

```json
{
  "success": true,
  "message": "File uploaded and processed successfully",
  "media": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "original_name": "photo.jpg",
    "stored_name": "123e4567-e89b-12d3-a456-426614174000_1704340800000000000.jpeg",
    "type": "image",
    "format": "jpeg",
    "size_bytes": 245680,
    "file_path": "./uploads/123e4567-e89b-12d3-a456-426614174000_1704340800000000000.jpeg",
    "uploaded_at": "2026-01-04T12:00:00Z",
    "width": 1920,
    "height": 1080
  }
}
```

**Response (Error - 400):**

```json
{
  "success": false,
  "error": "file size exceeds maximum limit of 200 MB (file size: 250.50 MB)"
}
```

### List All Media

**GET** `/media`

Retrieve all uploaded media files.

**Example with curl:**

```bash
curl http://localhost:8080/media
```

**Response:**

```json
{
  "total": 2,
  "media": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "original_name": "photo.jpg",
      "stored_name": "123e4567-e89b-12d3-a456-426614174000_1704340800000000000.jpeg",
      "type": "image",
      "format": "jpeg",
      "size_bytes": 245680,
      "file_path": "./uploads/123e4567-e89b-12d3-a456-426614174000_1704340800000000000.jpeg",
      "uploaded_at": "2026-01-04T12:00:00Z",
      "width": 1920,
      "height": 1080
    }
  ]
}
```

### Get Specific Media

**GET** `/media/{id}`

Retrieve details about a specific media file.

**Example with curl:**

```bash
curl http://localhost:8080/media/123e4567-e89b-12d3-a456-426614174000
```

### Download Media

**GET** `/media/{id}/download`

Download a specific media file.

**Example with curl:**

```bash
curl -O http://localhost:8080/media/123e4567-e89b-12d3-a456-426614174000/download
```

### Delete Media

**DELETE** `/media/{id}`

Delete a media file (removes from disk and database).

**Example with curl:**

```bash
curl -X DELETE http://localhost:8080/media/123e4567-e89b-12d3-a456-426614174000
```

**Response:**

```json
{
  "success": true,
  "message": "Media deleted successfully"
}
```

## Configuration

### Maximum File Size

The maximum file size is set to **200 MB**. This can be modified in [internal/services/media_service.go](internal/services/media_service.go):

```go
const MaxFileSize = 200 * 1024 * 1024
```

### Storage Location

Media files are stored in the `./uploads` directory, relative to the application root. This can be changed:

```go
const MediaStoragePath = "./uploads"
```

### Image Quality

JPEG optimization quality is set to **85%**. Adjust for your needs:

```go
const ImageQuality = 85
```

## Supported File Types

### Images

- JPEG/JPG
- PNG
- WebP
- GIF

### Documents

- PDF

All images are automatically converted to JPEG format during upload for optimal compression while maintaining good visual quality.

## File Structure

```
uploads/
├── [uuid]_[timestamp].jpeg    # Optimized JPEG image
├── [uuid]_[timestamp].jpeg    # Optimized JPEG image
└── [uuid]_[timestamp].pdf     # PDF document
```

## Error Handling

The service validates:

- **File Size**: Rejects files > 200 MB
- **File Type**: Only accepts JPEG, PNG, WebP, GIF, and PDF
- **Image Decoding**: Validates image integrity
- **Storage**: Checks filesystem permissions

## Usage Example

```go
// In your handler or service
mediaService := services.NewMediaService()

// Upload a file
file, fileHeader, _ := r.FormFile("file")
defer file.Close()

media, err := mediaService.UploadMedia(fileHeader)
if err != nil {
    // Handle error
    log.Printf("Upload failed: %v", err)
}

// Retrieve all media
allMedia := mediaService.GetAllMedia()

// Get specific media
media, err := mediaService.GetMedia(id)

// Delete media
err := mediaService.DeleteMedia(id)
```

## Notes

- All images are automatically optimized to JPEG format with 85% quality
- Image dimensions are preserved and stored in the metadata
- Files are stored with unique UUID-based names to prevent collisions
- The `./uploads` directory is created automatically if it doesn't exist
- Original filenames are preserved in the metadata for reference
