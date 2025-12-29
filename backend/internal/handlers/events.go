package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

// Event represents an event
type Event struct {
	ID           int64              `json:"id"`
	UUID         string             `json:"uuid"`
	State        string             `json:"state"`
	OnHome       bool               `json:"on_home"`
	Title        string             `json:"title"`
	Geolink      string             `json:"geolink,omitempty"`
	Type         string             `json:"type"`
	Location     string             `json:"location,omitempty"`
	StartDate    string             `json:"start_date"`
	EndDate      string             `json:"end_date"`
	ImageURL     string             `json:"imageUrl,omitempty"`
	TicketURL    string             `json:"ticket_purchase_url,omitempty"`
	Translations []EventTranslation `json:"translations,omitempty"`
}

type EventTranslation struct {
	LangCode    string `json:"lang_code"`
	Description string `json:"description"`
}

// GetPublicEvents returns all published events
func (h *Handler) GetPublicEvents(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")

	rows, err := h.db.Query(`
SELECT e.id, COALESCE(e.uuid, ''), e.state, e.on_home, e.title, e.geolink, e.type, e.location, 
       e.start_date, e.end_date, e.image_url, e.ticket_url
FROM events e
WHERE e.state = 'published'
ORDER BY e.start_date DESC
`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Event{})
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		var geolink, location, imageURL, ticketURL sql.NullString
		var onHome int

		if err := rows.Scan(
			&event.ID, &event.UUID, &event.State, &onHome, &event.Title,
			&geolink, &event.Type, &location, &event.StartDate,
			&event.EndDate, &imageURL, &ticketURL,
		); err != nil {
			continue
		}

		event.OnHome = onHome == 1
		if geolink.Valid {
			event.Geolink = geolink.String
		}
		if location.Valid {
			event.Location = location.String
		}
		if imageURL.Valid {
			event.ImageURL = imageURL.String
		}
		if ticketURL.Valid {
			event.TicketURL = ticketURL.String
		}

		// Get translations
		event.Translations = h.getEventTranslations(event.ID, lang)
		events = append(events, event)
	}

	respondJSON(w, http.StatusOK, events)
}

// GetHomeEvents returns events marked for homepage
func (h *Handler) GetHomeEvents(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")

	rows, err := h.db.Query(`
SELECT e.id, COALESCE(e.uuid, ''), e.state, e.on_home, e.title, e.geolink, e.type, e.location, 
       e.start_date, e.end_date, e.image_url, e.ticket_url
FROM events e
WHERE e.state = 'published' AND e.on_home = 1
ORDER BY e.start_date ASC
LIMIT 5
`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Event{})
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		var geolink, location, imageURL, ticketURL sql.NullString
		var onHome int

		if err := rows.Scan(
			&event.ID, &event.UUID, &event.State, &onHome, &event.Title,
			&geolink, &event.Type, &location, &event.StartDate,
			&event.EndDate, &imageURL, &ticketURL,
		); err != nil {
			continue
		}

		event.OnHome = onHome == 1
		if geolink.Valid {
			event.Geolink = geolink.String
		}
		if location.Valid {
			event.Location = location.String
		}
		if imageURL.Valid {
			event.ImageURL = imageURL.String
		}
		if ticketURL.Valid {
			event.TicketURL = ticketURL.String
		}

		event.Translations = h.getEventTranslations(event.ID, lang)
		events = append(events, event)
	}

	respondJSON(w, http.StatusOK, events)
}

// GetAdminEvents returns all events for admin
func (h *Handler) GetAdminEvents(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
SELECT id, COALESCE(uuid, ''), state, on_home, title, geolink, type, location, 
       start_date, end_date, image_url, ticket_url
FROM events
ORDER BY created_at DESC
`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Event{})
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		var geolink, location, imageURL, ticketURL sql.NullString
		var onHome int

		if err := rows.Scan(
			&event.ID, &event.UUID, &event.State, &onHome, &event.Title,
			&geolink, &event.Type, &location, &event.StartDate,
			&event.EndDate, &imageURL, &ticketURL,
		); err != nil {
			continue
		}

		event.OnHome = onHome == 1
		if geolink.Valid {
			event.Geolink = geolink.String
		}
		if location.Valid {
			event.Location = location.String
		}
		if imageURL.Valid {
			event.ImageURL = imageURL.String
		}
		if ticketURL.Valid {
			event.TicketURL = ticketURL.String
		}

		event.Translations = h.getEventTranslations(event.ID, "")
		events = append(events, event)
	}

	respondJSON(w, http.StatusOK, events)
}

// GetEventByID returns a single event (by numeric ID for admin)
func (h *Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var event Event
	var geolink, location, imageURL, ticketURL sql.NullString
	var onHome int

	err := h.db.QueryRow(`
SELECT id, COALESCE(uuid, ''), state, on_home, title, geolink, type, location, 
       start_date, end_date, image_url, ticket_url
FROM events WHERE id = ?
`, id).Scan(
		&event.ID, &event.UUID, &event.State, &onHome, &event.Title,
		&geolink, &event.Type, &location, &event.StartDate,
		&event.EndDate, &imageURL, &ticketURL,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Event not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	event.OnHome = onHome == 1
	if geolink.Valid {
		event.Geolink = geolink.String
	}
	if location.Valid {
		event.Location = location.String
	}
	if imageURL.Valid {
		event.ImageURL = imageURL.String
	}
	if ticketURL.Valid {
		event.TicketURL = ticketURL.String
	}

	event.Translations = h.getEventTranslations(event.ID, "")
	respondJSON(w, http.StatusOK, event)
}

// GetEventByUUID returns a single event by UUID (public)
func (h *Handler) GetEventByUUID(w http.ResponseWriter, r *http.Request) {
	eventUUID := chi.URLParam(r, "uuid")
	lang := r.URL.Query().Get("lang")
	preview := r.URL.Query().Get("preview") == "true"

	// Preview mode requires valid JWT authentication
	if preview {
		if !h.isAuthenticated(r) {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required for preview mode"})
			return
		}
	}

	var event Event
	var geolink, location, imageURL, ticketURL sql.NullString
	var onHome int

	// Allow preview of draft events if preview=true query param is set (already auth-checked)
	query := `
		SELECT id, COALESCE(uuid, ''), state, on_home, title, geolink, type, location, 
		       start_date, end_date, image_url, ticket_url
		FROM events WHERE uuid = ?`

	if !preview {
		query += ` AND state = 'published'`
	}

	err := h.db.QueryRow(query, eventUUID).Scan(
		&event.ID, &event.UUID, &event.State, &onHome, &event.Title,
		&geolink, &event.Type, &location, &event.StartDate,
		&event.EndDate, &imageURL, &ticketURL,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Event not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	event.OnHome = onHome == 1
	if geolink.Valid {
		event.Geolink = geolink.String
	}
	if location.Valid {
		event.Location = location.String
	}
	if imageURL.Valid {
		event.ImageURL = imageURL.String
	}
	if ticketURL.Valid {
		event.TicketURL = ticketURL.String
	}

	event.Translations = h.getEventTranslations(event.ID, lang)
	respondJSON(w, http.StatusOK, event)
}

// CreateEvent creates a new event
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	onHome := 0
	if event.OnHome {
		onHome = 1
	}

	// Generate UUID for new events
	event.UUID = uuid.New().String()

	result, err := h.db.Exec(`
INSERT INTO events (uuid, state, on_home, title, geolink, type, location, start_date, end_date, image_url, ticket_url)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, event.UUID, event.State, onHome, event.Title, event.Geolink, event.Type, event.Location,
		event.StartDate, event.EndDate, event.ImageURL, event.TicketURL)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create event"})
		return
	}

	event.ID, _ = result.LastInsertId()

	// Insert translations
	for _, t := range event.Translations {
		h.db.Exec(`
INSERT INTO event_translations (event_id, lang_code, description)
VALUES (?, ?, ?)
`, event.ID, t.LangCode, t.Description)
	}

	respondJSON(w, http.StatusCreated, event)
}

// UpdateEvent updates an existing event
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	onHome := 0
	if event.OnHome {
		onHome = 1
	}

	// Generate UUID if not exists
	var existingUUID sql.NullString
	h.db.QueryRow("SELECT uuid FROM events WHERE id = ?", id).Scan(&existingUUID)
	if !existingUUID.Valid || existingUUID.String == "" {
		event.UUID = uuid.New().String()
		_, err := h.db.Exec("UPDATE events SET uuid = ? WHERE id = ?", event.UUID, id)
		if err != nil {
			// Non-fatal, continue
		}
	} else {
		event.UUID = existingUUID.String
	}

	_, err := h.db.Exec(`
UPDATE events SET 
state = ?, on_home = ?, title = ?, geolink = ?, type = ?, location = ?,
start_date = ?, end_date = ?, image_url = ?, ticket_url = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`, event.State, onHome, event.Title, event.Geolink, event.Type, event.Location,
		event.StartDate, event.EndDate, event.ImageURL, event.TicketURL, id)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update event"})
		return
	}

	// Update translations
	for _, t := range event.Translations {
		h.db.Exec(`
INSERT OR REPLACE INTO event_translations (event_id, lang_code, description)
VALUES (?, ?, ?)
`, id, t.LangCode, t.Description)
	}

	idInt, _ := strconv.ParseInt(id, 10, 64)
	event.ID = idInt
	respondJSON(w, http.StatusOK, event)
}

// DeleteEvent deletes an event
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM events WHERE id = ?", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete event"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Event deleted"})
}

// Helper to get event translations
func (h *Handler) getEventTranslations(eventID int64, langFilter string) []EventTranslation {
	var rows *sql.Rows
	var err error

	if langFilter != "" {
		rows, err = h.db.Query(
			"SELECT lang_code, description FROM event_translations WHERE event_id = ? AND lang_code = ?",
			eventID, langFilter,
		)
	} else {
		rows, err = h.db.Query(
			"SELECT lang_code, description FROM event_translations WHERE event_id = ?",
			eventID,
		)
	}

	if err != nil {
		return []EventTranslation{}
	}
	defer rows.Close()

	translations := []EventTranslation{}
	for rows.Next() {
		var t EventTranslation
		if err := rows.Scan(&t.LangCode, &t.Description); err != nil {
			continue
		}
		translations = append(translations, t)
	}

	return translations
}

// GetEventQRCodes generates and returns a ZIP with QR codes for an event
func (h *Handler) GetEventQRCodes(w http.ResponseWriter, r *http.Request) {
	eventUUID := chi.URLParam(r, "uuid")

	// Verify event exists
	var exists bool
	err := h.db.QueryRow("SELECT 1 FROM events WHERE uuid = ?", eventUUID).Scan(&exists)
	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Event not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	// Construct event URL using frontend URL from config
	eventURL := fmt.Sprintf("%s/event/%s", h.cfg.FrontendURL, eventUUID)

	// Create ZIP buffer
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Generate 4 QR code variants
	variants := []struct {
		name   string
		fg     color.Color
		bg     color.Color
		transp bool
	}{
		{"qr-black-on-white.png", color.Black, color.White, false},
		{"qr-white-on-black.png", color.White, color.Black, false},
		{"qr-black-transparent.png", color.Black, color.Transparent, true},
		{"qr-white-transparent.png", color.White, color.Transparent, true},
	}

	for _, v := range variants {
		qr, err := qrcode.New(eventURL, qrcode.Medium)
		if err != nil {
			continue
		}
		qr.DisableBorder = true

		// Generate QR code image with custom colors
		qrImg := qr.Image(512)

		// Create new image with custom colors
		bounds := qrImg.Bounds()
		newImg := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				pixel := qrImg.At(x, y)
				r, g, b, _ := pixel.RGBA()

				// Check if pixel is dark (QR code module)
				isDark := (r+g+b)/3 < 32768

				if isDark {
					newImg.Set(x, y, v.fg)
				} else {
					newImg.Set(x, y, v.bg)
				}
			}
		}

		// Add to ZIP
		fileWriter, err := zipWriter.Create(v.name)
		if err != nil {
			continue
		}

		png.Encode(fileWriter, newImg)
	}

	zipWriter.Close()

	// Send ZIP file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"qr-codes-%s.zip\"", eventUUID))
	w.Write(buf.Bytes())
}
