package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/config"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/database"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/services/email"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Image magic byte signatures
var imageMagicBytes = map[string][]byte{
	".jpg":  {0xFF, 0xD8, 0xFF},
	".jpeg": {0xFF, 0xD8, 0xFF},
	".png":  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	".gif":  {0x47, 0x49, 0x46, 0x38},
	".webp": {0x52, 0x49, 0x46, 0x46}, // RIFF header (WEBP needs additional check)
}

// Handler holds dependencies for all HTTP handlers
type Handler struct {
	db           *database.DB
	cfg          *config.Config
	emailService *email.Service
}

// New creates a new Handler with all dependencies
func New(db *database.DB, cfg *config.Config, emailService *email.Service) *Handler {
	return &Handler{
		db:           db,
		cfg:          cfg,
		emailService: emailService,
	}
}

// UploadImage handles image file uploads
func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Max 10MB file
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("image")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "No image file provided"})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".gif" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type. Allowed: jpg, jpeg, png, webp, gif"})
		return
	}

	// Read file content for magic byte validation
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read file"})
		return
	}

	// Validate magic bytes
	if !validateImageMagicBytes(fileContent, ext) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "File content does not match declared type"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create images directory if it doesn't exist
	imagesDir := filepath.Join(h.cfg.DataDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create images directory"})
		return
	}

	// Save file
	dst, err := os.Create(filepath.Join(imagesDir, filename))
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save image"})
		return
	}
	defer dst.Close()

	if _, err := dst.Write(fileContent); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save image"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"filename": filename})
}

// validateImageMagicBytes checks if file content matches the expected image type
func validateImageMagicBytes(content []byte, ext string) bool {
	if len(content) < 12 {
		return false
	}

	magic, exists := imageMagicBytes[ext]
	if !exists {
		return false
	}

	// Check magic bytes
	if !bytes.HasPrefix(content, magic) {
		return false
	}

	// Additional check for WEBP (RIFF header + WEBP signature)
	if ext == ".webp" {
		// WEBP format: RIFF....WEBP
		if len(content) < 12 {
			return false
		}
		if string(content[8:12]) != "WEBP" {
			return false
		}
	}

	return true
}

// ServeImage serves uploaded images
func (h *Handler) ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/images/")
	if filename == "" {
		http.NotFound(w, r)
		return
	}

	// Sanitize filename to prevent directory traversal
	filename = filepath.Base(filename)
	imagePath := filepath.Join(h.cfg.DataDir, "images", filename)

	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Set cache headers
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	http.ServeFile(w, r, imagePath)
}

// isAuthenticated checks if request has a valid JWT token (for optional auth on public endpoints)
func (h *Handler) isAuthenticated(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.cfg.JWT.Secret), nil
	})

	return err == nil && token.Valid
}
