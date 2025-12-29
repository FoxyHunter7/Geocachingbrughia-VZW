package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Language represents a supported language
type Language struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	FlagURL string `json:"flag_url,omitempty"`
	Active  bool   `json:"active,omitempty"`
}

// StaticContent represents a translation entry
type StaticContent struct {
	Property string                   `json:"property"`
	Contents []StaticContentTranslation `json:"contents"`
}

type StaticContentTranslation struct {
	LangCode string `json:"lang_code"`
	Content  string `json:"content"`
}

// Social represents a social media link
type Social struct {
	ID        int64  `json:"id"`
	Platform  string `json:"platform"`
	URL       string `json:"url"`
	Icon      string `json:"icon,omitempty"`
	Active    bool   `json:"active,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
}

// GetLanguages returns all active languages (public)
func (h *Handler) GetLanguages(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT code, name, flag_url FROM languages WHERE active = 1")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Language{})
		return
	}
	defer rows.Close()

	languages := []Language{}
	for rows.Next() {
		var lang Language
		var flagURL *string
		if err := rows.Scan(&lang.Code, &lang.Name, &flagURL); err != nil {
			continue
		}
		if flagURL != nil {
			lang.FlagURL = *flagURL
		}
		languages = append(languages, lang)
	}

	// Set ETag based on data
	etag := generateETag(languages)
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("ETag", etag)

	respondJSON(w, http.StatusOK, languages)
}

// GetStaticContent returns all static content/translations (public)
func (h *Handler) GetStaticContent(w http.ResponseWriter, r *http.Request) {
	// Get the latest update time for ETag
	var lastUpdate time.Time
	h.db.QueryRow("SELECT MAX(updated_at) FROM static_content").Scan(&lastUpdate)
	
	etag := fmt.Sprintf(`"%x"`, md5.Sum([]byte(lastUpdate.String())))
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("ETag", etag)

	rows, err := h.db.Query("SELECT property, lang_code, content FROM static_content ORDER BY property")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []StaticContent{})
		return
	}
	defer rows.Close()

	// Group by property
	contentMap := make(map[string]*StaticContent)
	for rows.Next() {
		var property, langCode, content string
		if err := rows.Scan(&property, &langCode, &content); err != nil {
			continue
		}

		if _, exists := contentMap[property]; !exists {
			contentMap[property] = &StaticContent{
				Property: property,
				Contents: []StaticContentTranslation{},
			}
		}
		contentMap[property].Contents = append(contentMap[property].Contents, StaticContentTranslation{
			LangCode: langCode,
			Content:  content,
		})
	}

	// Convert map to slice
	result := make([]StaticContent, 0, len(contentMap))
	for _, v := range contentMap {
		result = append(result, *v)
	}

	respondJSON(w, http.StatusOK, result)
}

// GetSocials returns all active social media links (public)
func (h *Handler) GetSocials(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, platform, url, icon FROM socials WHERE active = 1 ORDER BY sort_order")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Social{})
		return
	}
	defer rows.Close()

	socials := []Social{}
	for rows.Next() {
		var social Social
		var icon *string
		if err := rows.Scan(&social.ID, &social.Platform, &social.URL, &icon); err != nil {
			continue
		}
		if icon != nil {
			social.Icon = *icon
		}
		socials = append(socials, social)
	}

	etag := generateETag(socials)
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("ETag", etag)

	respondJSON(w, http.StatusOK, socials)
}

// Admin handlers for languages
func (h *Handler) GetAdminLanguages(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT code, name, flag_url, active FROM languages")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Language{})
		return
	}
	defer rows.Close()

	languages := []Language{}
	for rows.Next() {
		var lang Language
		var flagURL *string
		var active int
		if err := rows.Scan(&lang.Code, &lang.Name, &flagURL, &active); err != nil {
			continue
		}
		if flagURL != nil {
			lang.FlagURL = *flagURL
		}
		lang.Active = active == 1
		languages = append(languages, lang)
	}

	respondJSON(w, http.StatusOK, languages)
}

func (h *Handler) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	var lang Language
	if err := json.NewDecoder(r.Body).Decode(&lang); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	active := 0
	if lang.Active {
		active = 1
	}

	_, err := h.db.Exec(
		"INSERT INTO languages (code, name, flag_url, active) VALUES (?, ?, ?, ?)",
		lang.Code, lang.Name, lang.FlagURL, active,
	)
	if err != nil {
		respondJSON(w, http.StatusConflict, map[string]string{"error": "Language code already exists"})
		return
	}

	respondJSON(w, http.StatusCreated, lang)
}

func (h *Handler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		// Fallback for older chi versions
		code = getPathParam(r, "code")
	}

	var lang Language
	if err := json.NewDecoder(r.Body).Decode(&lang); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	active := 0
	if lang.Active {
		active = 1
	}

	_, err := h.db.Exec(
		"UPDATE languages SET name = ?, flag_url = ?, active = ?, updated_at = CURRENT_TIMESTAMP WHERE code = ?",
		lang.Name, lang.FlagURL, active, code,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update language"})
		return
	}

	respondJSON(w, http.StatusOK, lang)
}

func (h *Handler) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		code = getPathParam(r, "code")
	}

	_, err := h.db.Exec("DELETE FROM languages WHERE code = ?", code)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete language"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Language deleted"})
}

// Admin handlers for static content
func (h *Handler) GetAdminStaticContent(w http.ResponseWriter, r *http.Request) {
	// Same as public but might include additional metadata in future
	h.GetStaticContent(w, r)
}

func (h *Handler) CreateStaticContent(w http.ResponseWriter, r *http.Request) {
	var content struct {
		Property string `json:"property"`
		LangCode string `json:"lang_code"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err := h.db.Exec(
		"INSERT INTO static_content (property, lang_code, content) VALUES (?, ?, ?)",
		content.Property, content.LangCode, content.Content,
	)
	if err != nil {
		respondJSON(w, http.StatusConflict, map[string]string{"error": "Content already exists for this property/language"})
		return
	}

	respondJSON(w, http.StatusCreated, content)
}

func (h *Handler) UpdateStaticContent(w http.ResponseWriter, r *http.Request) {
	property := r.PathValue("property")
	if property == "" {
		property = getPathParam(r, "property")
	}

	var content struct {
		LangCode string `json:"lang_code"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err := h.db.Exec(
		"UPDATE static_content SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE property = ? AND lang_code = ?",
		content.Content, property, content.LangCode,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update content"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Content updated"})
}

func (h *Handler) DeleteStaticContent(w http.ResponseWriter, r *http.Request) {
	property := r.PathValue("property")
	if property == "" {
		property = getPathParam(r, "property")
	}

	// Delete all translations for this property
	_, err := h.db.Exec("DELETE FROM static_content WHERE property = ?", property)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete content"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Content deleted"})
}

// Admin handlers for socials
func (h *Handler) GetAdminSocials(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, platform, url, icon, active, sort_order FROM socials ORDER BY sort_order")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Social{})
		return
	}
	defer rows.Close()

	socials := []Social{}
	for rows.Next() {
		var social Social
		var icon *string
		var active int
		if err := rows.Scan(&social.ID, &social.Platform, &social.URL, &icon, &active, &social.SortOrder); err != nil {
			continue
		}
		if icon != nil {
			social.Icon = *icon
		}
		social.Active = active == 1
		socials = append(socials, social)
	}

	respondJSON(w, http.StatusOK, socials)
}

func (h *Handler) CreateSocial(w http.ResponseWriter, r *http.Request) {
	var social Social
	if err := json.NewDecoder(r.Body).Decode(&social); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	active := 0
	if social.Active {
		active = 1
	}

	result, err := h.db.Exec(
		"INSERT INTO socials (platform, url, icon, active, sort_order) VALUES (?, ?, ?, ?, ?)",
		social.Platform, social.URL, social.Icon, active, social.SortOrder,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create social"})
		return
	}

	social.ID, _ = result.LastInsertId()
	respondJSON(w, http.StatusCreated, social)
}

func (h *Handler) UpdateSocial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		id = getPathParam(r, "id")
	}

	var social Social
	if err := json.NewDecoder(r.Body).Decode(&social); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	active := 0
	if social.Active {
		active = 1
	}

	_, err := h.db.Exec(
		"UPDATE socials SET platform = ?, url = ?, icon = ?, active = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		social.Platform, social.URL, social.Icon, active, social.SortOrder, id,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update social"})
		return
	}

	respondJSON(w, http.StatusOK, social)
}

func (h *Handler) DeleteSocial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		id = getPathParam(r, "id")
	}

	_, err := h.db.Exec("DELETE FROM socials WHERE id = ?", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete social"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Social deleted"})
}

// Helper to generate ETag from data
func generateETag(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return fmt.Sprintf(`"%x"`, md5.Sum(jsonData))
}

// Helper to get path parameter (chi router)
func getPathParam(r *http.Request, name string) string {
	// Chi stores URL params in context
	if ctx := r.Context(); ctx != nil {
		if val := ctx.Value(name); val != nil {
			return val.(string)
		}
	}
	return ""
}
