package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Geocache represents a geocache
type Geocache struct {
	ID         int64   `json:"id"`
	GCCode     string  `json:"gc_code"`
	Name       string  `json:"name"`
	Title      string  `json:"title"`                 // Alias for name for frontend compatibility
	Geolink    string  `json:"geolink"`               // Full geocaching.com URL
	Type       string  `json:"type"`                  // Cache type (traditional, mystery, multi, etc.)
	PlacedDate string  `json:"placed_date,omitempty"` // Date the cache was placed
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Difficulty float64 `json:"difficulty,omitempty"`
	Terrain    float64 `json:"terrain,omitempty"`
	Size       string  `json:"size,omitempty"`
	Status     string  `json:"status"`
}

// GetPublicGeocaches returns all active geocaches
func (h *Handler) GetPublicGeocaches(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, gc_code, name, latitude, longitude, difficulty, terrain, size, type, placed_date, status
		FROM geocaches
		WHERE status = 'active'
		ORDER BY name
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Geocache{})
		return
	}
	defer rows.Close()

	geocaches := h.scanGeocaches(rows)
	respondJSON(w, http.StatusOK, geocaches)
}

// GetAdminGeocaches returns all geocaches for admin
func (h *Handler) GetAdminGeocaches(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, gc_code, name, latitude, longitude, difficulty, terrain, size, type, placed_date, status
		FROM geocaches
		ORDER BY created_at DESC
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Geocache{})
		return
	}
	defer rows.Close()

	geocaches := h.scanGeocaches(rows)
	respondJSON(w, http.StatusOK, geocaches)
}

// GetGeocacheByID returns a single geocache
func (h *Handler) GetGeocacheByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var gc Geocache
	var lat, lng, diff, terr sql.NullFloat64
	var size, gcType, placedDate sql.NullString

	err := h.db.QueryRow(`
		SELECT id, gc_code, name, latitude, longitude, difficulty, terrain, size, type, placed_date, status
		FROM geocaches WHERE id = ?
	`, id).Scan(&gc.ID, &gc.GCCode, &gc.Name, &lat, &lng, &diff, &terr, &size, &gcType, &placedDate, &gc.Status)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Geocache not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	// Set title alias
	gc.Title = gc.Name
	// gc_code now stores the full geolink URL
	gc.Geolink = gc.GCCode

	if lat.Valid {
		gc.Latitude = lat.Float64
	}
	if lng.Valid {
		gc.Longitude = lng.Float64
	}
	if diff.Valid {
		gc.Difficulty = diff.Float64
	}
	if terr.Valid {
		gc.Terrain = terr.Float64
	}
	if size.Valid {
		gc.Size = size.String
	}
	if gcType.Valid {
		gc.Type = gcType.String
	} else {
		gc.Type = "traditional"
	}
	if placedDate.Valid {
		gc.PlacedDate = placedDate.String
	}

	respondJSON(w, http.StatusOK, gc)
}

// CreateGeocache creates a new geocache
func (h *Handler) CreateGeocache(w http.ResponseWriter, r *http.Request) {
	var gc Geocache
	if err := json.NewDecoder(r.Body).Decode(&gc); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if gc.Status == "" {
		gc.Status = "active"
	}
	if gc.Type == "" {
		gc.Type = "traditional"
	}

	// Store geolink in gc_code column
	result, err := h.db.Exec(`
		INSERT INTO geocaches (gc_code, name, latitude, longitude, difficulty, terrain, size, type, placed_date, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, gc.Geolink, gc.Name, gc.Latitude, gc.Longitude, gc.Difficulty, gc.Terrain, gc.Size, gc.Type, gc.PlacedDate, gc.Status)

	if err != nil {
		respondJSON(w, http.StatusConflict, map[string]string{"error": "Link already exists"})
		return
	}

	gc.ID, _ = result.LastInsertId()
	gc.Title = gc.Name
	respondJSON(w, http.StatusCreated, gc)
}

// UpdateGeocache updates an existing geocache
func (h *Handler) UpdateGeocache(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var gc Geocache
	if err := json.NewDecoder(r.Body).Decode(&gc); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if gc.Type == "" {
		gc.Type = "traditional"
	}

	// Store geolink in gc_code column
	_, err := h.db.Exec(`
		UPDATE geocaches SET 
			gc_code = ?, name = ?, latitude = ?, longitude = ?, 
			difficulty = ?, terrain = ?, size = ?, type = ?, placed_date = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, gc.Geolink, gc.Name, gc.Latitude, gc.Longitude, gc.Difficulty, gc.Terrain, gc.Size, gc.Type, gc.PlacedDate, gc.Status, id)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update geocache"})
		return
	}

	idInt, _ := strconv.ParseInt(id, 10, 64)
	gc.ID = idInt
	gc.Title = gc.Name
	respondJSON(w, http.StatusOK, gc)
}

// DeleteGeocache deletes a geocache
func (h *Handler) DeleteGeocache(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM geocaches WHERE id = ?", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete geocache"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Geocache deleted"})
}

// Helper to scan geocaches from rows
func (h *Handler) scanGeocaches(rows *sql.Rows) []Geocache {
	geocaches := []Geocache{}
	for rows.Next() {
		var gc Geocache
		var lat, lng, diff, terr sql.NullFloat64
		var size, gcType, placedDate sql.NullString

		if err := rows.Scan(&gc.ID, &gc.GCCode, &gc.Name, &lat, &lng, &diff, &terr, &size, &gcType, &placedDate, &gc.Status); err != nil {
			continue
		}

		// Set title alias
		gc.Title = gc.Name
		// gc_code now stores the full geolink URL
		gc.Geolink = gc.GCCode

		if lat.Valid {
			gc.Latitude = lat.Float64
		}
		if lng.Valid {
			gc.Longitude = lng.Float64
		}
		if diff.Valid {
			gc.Difficulty = diff.Float64
		}
		if terr.Valid {
			gc.Terrain = terr.Float64
		}
		if size.Valid {
			gc.Size = size.String
		}
		if gcType.Valid {
			gc.Type = gcType.String
		} else {
			gc.Type = "traditional" // Default type
		}
		if placedDate.Valid {
			gc.PlacedDate = placedDate.String
		}

		geocaches = append(geocaches, gc)
	}
	return geocaches
}
