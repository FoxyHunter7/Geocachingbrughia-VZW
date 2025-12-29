package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Message represents a site message/announcement
type Message struct {
	ID           int64                `json:"id"`
	State        string               `json:"state"`
	Priority     int                  `json:"priority"`
	Translations []MessageTranslation `json:"translations,omitempty"`
}

type MessageTranslation struct {
	LangCode string `json:"lang_code"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// GetPublicMessages returns all published messages
func (h *Handler) GetPublicMessages(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")

	rows, err := h.db.Query(`
		SELECT id, state, priority
		FROM messages
		WHERE state = 'published'
		ORDER BY priority DESC, created_at DESC
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Message{})
		return
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.State, &msg.Priority); err != nil {
			continue
		}
		msg.Translations = h.getMessageTranslations(msg.ID, lang)
		messages = append(messages, msg)
	}

	respondJSON(w, http.StatusOK, messages)
}

// GetAdminMessages returns all messages for admin
func (h *Handler) GetAdminMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, state, priority
		FROM messages
		ORDER BY created_at DESC
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, []Message{})
		return
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.State, &msg.Priority); err != nil {
			continue
		}
		msg.Translations = h.getMessageTranslations(msg.ID, "")
		messages = append(messages, msg)
	}

	respondJSON(w, http.StatusOK, messages)
}

// GetMessageByID returns a single message
func (h *Handler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var msg Message
	err := h.db.QueryRow(`
		SELECT id, state, priority
		FROM messages WHERE id = ?
	`, id).Scan(&msg.ID, &msg.State, &msg.Priority)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Message not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	msg.Translations = h.getMessageTranslations(msg.ID, "")
	respondJSON(w, http.StatusOK, msg)
}

// CreateMessage creates a new message
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if msg.State == "" {
		msg.State = "draft"
	}

	result, err := h.db.Exec(`
		INSERT INTO messages (state, priority)
		VALUES (?, ?)
	`, msg.State, msg.Priority)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create message"})
		return
	}

	msg.ID, _ = result.LastInsertId()

	// Insert translations
	for _, t := range msg.Translations {
		h.db.Exec(`
			INSERT INTO message_translations (message_id, lang_code, title, content)
			VALUES (?, ?, ?, ?)
		`, msg.ID, t.LangCode, t.Title, t.Content)
	}

	respondJSON(w, http.StatusCreated, msg)
}

// UpdateMessage updates an existing message
func (h *Handler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE messages SET state = ?, priority = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, msg.State, msg.Priority, id)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update message"})
		return
	}

	// Update translations
	for _, t := range msg.Translations {
		h.db.Exec(`
			INSERT OR REPLACE INTO message_translations (message_id, lang_code, title, content)
			VALUES (?, ?, ?, ?)
		`, id, t.LangCode, t.Title, t.Content)
	}

	idInt, _ := strconv.ParseInt(id, 10, 64)
	msg.ID = idInt
	respondJSON(w, http.StatusOK, msg)
}

// DeleteMessage deletes a message
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete message"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Message deleted"})
}

// Helper to get message translations
func (h *Handler) getMessageTranslations(messageID int64, langFilter string) []MessageTranslation {
	var rows *sql.Rows
	var err error

	if langFilter != "" {
		rows, err = h.db.Query(
			"SELECT lang_code, title, content FROM message_translations WHERE message_id = ? AND lang_code = ?",
			messageID, langFilter,
		)
	} else {
		rows, err = h.db.Query(
			"SELECT lang_code, title, content FROM message_translations WHERE message_id = ?",
			messageID,
		)
	}

	if err != nil {
		return []MessageTranslation{}
	}
	defer rows.Close()

	translations := []MessageTranslation{}
	for rows.Next() {
		var t MessageTranslation
		if err := rows.Scan(&t.LangCode, &t.Title, &t.Content); err != nil {
			continue
		}
		translations = append(translations, t)
	}

	return translations
}
