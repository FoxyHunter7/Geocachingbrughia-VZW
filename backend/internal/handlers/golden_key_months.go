package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GoldenKeyMonth represents a monthly golden key entry
type GoldenKeyMonth struct {
	ID          int64      `json:"id"`
	MonthNumber int        `json:"month_number"`
	MonthName   string     `json:"month_name"`
	LiveDate    time.Time  `json:"live_date"`
	FoundDate   *time.Time `json:"found_date,omitempty"`
	State       string     `json:"state"` // computed: locked, active, found
	IsFound     bool       `json:"is_found"`
	FinderName  string     `json:"finder_name,omitempty"`
	FinderImage string     `json:"finder_image,omitempty"`
}

// GoldenKeyHint is a single hint for a month
type GoldenKeyHint struct {
	ID        int64  `json:"id"`
	MonthID   int64  `json:"month_id"`
	SortOrder int    `json:"sort_order"`
	Content   string `json:"content"`
	ImageURL  string `json:"image_url,omitempty"`
}

// GoldenKeyMonthDetail includes hints
type GoldenKeyMonthDetail struct {
	GoldenKeyMonth
	Hints []GoldenKeyHint `json:"hints"`
}

func computeMonthState(liveDate time.Time, isFound int) string {
	if isFound == 1 {
		return "found"
	}
	if liveDate.IsZero() {
		return "locked"
	}
	if time.Now().UTC().After(liveDate) {
		return "active"
	}
	return "locked"
}

func scanMonth(rows interface {
	Scan(dest ...any) error
}) (GoldenKeyMonth, time.Time, int, error) {
	var m GoldenKeyMonth
	var liveDateStr string
	var isFound int
	var finderName, finderImage, foundDate sql.NullString
	err := rows.Scan(&m.ID, &m.MonthNumber, &m.MonthName, &liveDateStr, &isFound, &finderName, &finderImage, &foundDate)
	if finderName.Valid {
		m.FinderName = finderName.String
	}
	if finderImage.Valid {
		m.FinderImage = finderImage.String
	}
	liveDate, _ := parseFlexibleTime(liveDateStr)
	m.LiveDate = liveDate
	if foundDate.Valid && foundDate.String != "" {
		if fd, perr := parseFlexibleTime(foundDate.String); perr == nil {
			m.FoundDate = &fd
		}
	}
	return m, liveDate, isFound, err
}

// --- Public endpoints ---

// GetGoldenKeyMonths returns all 12 months with computed state (public).
func (h *Handler) GetGoldenKeyMonths(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, month_number, month_name, live_date, is_found, finder_name, finder_image, found_date
		FROM golden_key_months ORDER BY month_number
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []GoldenKeyMonth{})
		return
	}
	defer rows.Close()

	months := []GoldenKeyMonth{}
	for rows.Next() {
		m, liveDate, isFound, err := scanMonth(rows)
		if err != nil {
			continue
		}
		m.IsFound = isFound == 1
		m.State = computeMonthState(liveDate, isFound)
		if m.State != "found" {
			m.FinderName = ""
			m.FinderImage = ""
		}
		months = append(months, m)
	}
	respondJSON(w, http.StatusOK, months)
}

// GetGoldenKeyMonthByID returns a specific month with hints (public).
// Returns 403 if the month is still locked.
func (h *Handler) GetGoldenKeyMonthByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	m, liveDate, isFound, err := scanMonth(h.db.QueryRow(`
		SELECT id, month_number, month_name, live_date, is_found, finder_name, finder_image, found_date
		FROM golden_key_months WHERE id = ?
	`, id))
	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Month not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch month"})
		return
	}

	m.IsFound = isFound == 1
	m.State = computeMonthState(liveDate, isFound)

	if m.State == "locked" {
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "locked", "state": "locked"})
		return
	}
	if m.State != "found" {
		m.FinderName = ""
		m.FinderImage = ""
	}

	hints, err := fetchHints(h, id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch hints"})
		return
	}

	respondJSON(w, http.StatusOK, GoldenKeyMonthDetail{GoldenKeyMonth: m, Hints: hints})
}

// --- Admin endpoints ---

// GetAdminGoldenKeyMonths returns all months with full data for admin.
func (h *Handler) GetAdminGoldenKeyMonths(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, month_number, month_name, live_date, is_found, finder_name, finder_image, found_date
		FROM golden_key_months ORDER BY month_number
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []GoldenKeyMonth{})
		return
	}
	defer rows.Close()

	months := []GoldenKeyMonth{}
	for rows.Next() {
		m, liveDate, isFound, err := scanMonth(rows)
		if err != nil {
			continue
		}
		m.IsFound = isFound == 1
		m.State = computeMonthState(liveDate, isFound)
		months = append(months, m)
	}
	respondJSON(w, http.StatusOK, months)
}

// GetAdminGoldenKeyMonthByID returns a month with hints for admin.
func (h *Handler) GetAdminGoldenKeyMonthByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	m, liveDate, isFound, err := scanMonth(h.db.QueryRow(`
		SELECT id, month_number, month_name, live_date, is_found, finder_name, finder_image, found_date
		FROM golden_key_months WHERE id = ?
	`, id))
	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Month not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch month"})
		return
	}

	m.IsFound = isFound == 1
	m.State = computeMonthState(liveDate, isFound)

	hints, err := fetchHints(h, id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch hints"})
		return
	}
	respondJSON(w, http.StatusOK, GoldenKeyMonthDetail{GoldenKeyMonth: m, Hints: hints})
}

type updateMonthRequest struct {
	LiveDate    string `json:"live_date"`
	FoundDate   string `json:"found_date"`
	IsFound     bool   `json:"is_found"`
	FinderName  string `json:"finder_name"`
	FinderImage string `json:"finder_image"`
}

// UpdateGoldenKeyMonth updates month settings (admin).
func (h *Handler) UpdateGoldenKeyMonth(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var req updateMonthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	liveDate, err := time.Parse(time.RFC3339, req.LiveDate)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid live_date, expected RFC3339"})
		return
	}

	isFound := 0
	if req.IsFound {
		isFound = 1
	}

	var foundDateVal interface{}
	if req.FoundDate != "" {
		fd, ferr := time.Parse(time.RFC3339, req.FoundDate)
		if ferr != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid found_date, expected RFC3339"})
			return
		}
		foundDateVal = fd.UTC().Format("2006-01-02 15:04:05")
	}

	_, err = h.db.Exec(`
		UPDATE golden_key_months
		SET live_date = ?, is_found = ?, finder_name = ?, finder_image = ?, found_date = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, liveDate.UTC().Format("2006-01-02 15:04:05"), isFound, req.FinderName, req.FinderImage, foundDateVal, id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update month"})
		return
	}

	// Return the updated month with hints
	h.GetAdminGoldenKeyMonthByID(w, r)
}

// --- Hint endpoints ---

type hintRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
}

// AddGoldenKeyHint adds a hint to a month (admin).
func (h *Handler) AddGoldenKeyHint(w http.ResponseWriter, r *http.Request) {
	monthID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid month ID"})
		return
	}

	var req hintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	var maxOrder int
	h.db.QueryRow(`SELECT COALESCE(MAX(sort_order), -1) FROM golden_key_hints WHERE month_id = ?`, monthID).Scan(&maxOrder)

	result, err := h.db.Exec(`
		INSERT INTO golden_key_hints (month_id, sort_order, content, image_url)
		VALUES (?, ?, ?, ?)
	`, monthID, maxOrder+1, req.Content, nullableString(req.ImageURL))
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add hint"})
		return
	}

	hintID, _ := result.LastInsertId()
	respondJSON(w, http.StatusCreated, GoldenKeyHint{
		ID:        hintID,
		MonthID:   monthID,
		SortOrder: maxOrder + 1,
		Content:   req.Content,
		ImageURL:  req.ImageURL,
	})
}

// UpdateGoldenKeyHint updates hint content (admin).
func (h *Handler) UpdateGoldenKeyHint(w http.ResponseWriter, r *http.Request) {
	hintID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid hint ID"})
		return
	}

	var req hintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err = h.db.Exec(`
		UPDATE golden_key_hints SET content = ?, image_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, req.Content, nullableString(req.ImageURL), hintID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update hint"})
		return
	}

	var hint GoldenKeyHint
	var imgURL sql.NullString
	h.db.QueryRow(`SELECT id, month_id, sort_order, content, image_url FROM golden_key_hints WHERE id = ?`, hintID).
		Scan(&hint.ID, &hint.MonthID, &hint.SortOrder, &hint.Content, &imgURL)
	if imgURL.Valid {
		hint.ImageURL = imgURL.String
	}
	respondJSON(w, http.StatusOK, hint)
}

// DeleteGoldenKeyHint deletes a hint (admin).
func (h *Handler) DeleteGoldenKeyHint(w http.ResponseWriter, r *http.Request) {
	hintID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid hint ID"})
		return
	}

	if _, err := h.db.Exec(`DELETE FROM golden_key_hints WHERE id = ?`, hintID); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete hint"})
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// --- helpers ---

func fetchHints(h *Handler, monthID int64) ([]GoldenKeyHint, error) {
	rows, err := h.db.Query(`
		SELECT id, month_id, sort_order, content, image_url
		FROM golden_key_hints WHERE month_id = ? ORDER BY sort_order, id
	`, monthID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hints := []GoldenKeyHint{}
	for rows.Next() {
		var hint GoldenKeyHint
		var imgURL sql.NullString
		if err := rows.Scan(&hint.ID, &hint.MonthID, &hint.SortOrder, &hint.Content, &imgURL); err != nil {
			continue
		}
		if imgURL.Valid {
			hint.ImageURL = imgURL.String
		}
		hints = append(hints, hint)
	}
	return hints, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
