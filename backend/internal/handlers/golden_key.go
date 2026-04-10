package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type goldenKeySettings struct {
	ActivationTime time.Time         `json:"activation_time"`
	IsActive       bool              `json:"is_active"`
	BannerText     map[string]string `json:"banner_text"`
	Rules          map[string]string `json:"rules"`
}

type updateGoldenKeyRequest struct {
	ActivationTime string            `json:"activation_time"`
	BannerText     map[string]string `json:"banner_text"`
	Rules          map[string]string `json:"rules"`
}

// GetGoldenKeySettings returns the golden key activation time and active status.
// Public endpoint — no authentication required.
func (h *Handler) GetGoldenKeySettings(w http.ResponseWriter, r *http.Request) {
	var activationTimeStr string
	var bannerTextJSON string
	var rulesJSON string

	err := h.db.QueryRow(`SELECT activation_time, banner_text, rules FROM golden_key_settings WHERE id = 1`).Scan(&activationTimeStr, &bannerTextJSON, &rulesJSON)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch golden key settings"})
		return
	}

	activationTime, err := parseFlexibleTime(activationTimeStr)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse activation time"})
		return
	}

	bannerTextMap := map[string]string{}
	if bannerTextJSON != "" {
		_ = json.Unmarshal([]byte(bannerTextJSON), &bannerTextMap)
	}

	rulesMap := map[string]string{}
	if rulesJSON != "" {
		_ = json.Unmarshal([]byte(rulesJSON), &rulesMap)
	}

	respondJSON(w, http.StatusOK, goldenKeySettings{
		ActivationTime: activationTime,
		IsActive:       time.Now().UTC().After(activationTime),
		BannerText:     bannerTextMap,
		Rules:          rulesMap,
	})
}

// UpdateGoldenKeySettings updates the golden key activation time.
// Protected admin endpoint.
func (h *Handler) UpdateGoldenKeySettings(w http.ResponseWriter, r *http.Request) {
	var req updateGoldenKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	activationTime, err := time.Parse(time.RFC3339, req.ActivationTime)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid activation_time format, expected RFC3339 (e.g. 2026-04-12T10:12:00Z)"})
		return
	}

	utcStr := activationTime.UTC().Format("2006-01-02 15:04:05")

	bannerTextBytes, _ := json.Marshal(req.BannerText)
	rulesBytes, _ := json.Marshal(req.Rules)

	_, err = h.db.Exec(
		`UPDATE golden_key_settings SET activation_time = ?, banner_text = ?, rules = ?, updated_at = CURRENT_TIMESTAMP WHERE id = 1`,
		utcStr, string(bannerTextBytes), string(rulesBytes),
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update golden key settings"})
		return
	}

	h.GetGoldenKeySettings(w, r)
}

// parseFlexibleTime tries multiple SQLite datetime formats before giving up.
func parseFlexibleTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, &time.ParseError{Value: s}
}
