package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// ContactSubmission represents a contact form submission
type ContactSubmission struct {
	ID                 int64         `json:"id"`
	Email              string        `json:"email"`
	Subject            string        `json:"subject"`
	Message            string        `json:"message"`
	Status             string        `json:"status"`
	AssignedTo         *int64        `json:"assigned_to,omitempty"`
	AssignedToName     string        `json:"assigned_to_name,omitempty"`
	LastReminderSentAt *string       `json:"last_reminder_sent_at,omitempty"`
	CreatedAt          string        `json:"created_at"`
	UpdatedAt          string        `json:"updated_at"`
	Notes              []ContactNote `json:"notes,omitempty"`
}

type ContactNote struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

type ContactFormRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// SubmitContactForm handles public contact form submissions
func (h *Handler) SubmitContactForm(w http.ResponseWriter, r *http.Request) {
	var req ContactFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate input
	errors := make(map[string][]string)
	if req.Email == "" {
		errors["email"] = append(errors["email"], "Email is required")
	}
	if req.Subject == "" {
		errors["subject"] = append(errors["subject"], "Subject is required")
	}
	if req.Message == "" {
		errors["message"] = append(errors["message"], "Message is required")
	}

	if len(errors) > 0 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"errors": errors,
		})
		return
	}

	// Insert into database
	result, err := h.db.Exec(`
		INSERT INTO contact_submissions (email, subject, message, status)
		VALUES (?, ?, ?, 'new')
	`, req.Email, req.Subject, req.Message)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to submit contact form"})
		return
	}

	submissionID, _ := result.LastInsertId()

	// Send notification email (async, don't block the response)
	go h.emailService.SendNewContactNotification(req.Email, req.Subject, req.Message, submissionID)

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"status": true,
	})
}

// GetContactSubmissions returns all contact submissions for admin
func (h *Handler) GetContactSubmissions(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 15
	offset := (page - 1) * perPage

	var rows *sql.Rows
	var err error
	var totalCount int

	if status != "" {
		h.db.QueryRow("SELECT COUNT(*) FROM contact_submissions WHERE status = ?", status).Scan(&totalCount)
		rows, err = h.db.Query(`
			SELECT cs.id, cs.email, cs.subject, cs.message, cs.status, 
			       cs.assigned_to, u.name, cs.last_reminder_sent_at, cs.created_at, cs.updated_at
			FROM contact_submissions cs
			LEFT JOIN users u ON cs.assigned_to = u.id
			WHERE cs.status = ?
			ORDER BY cs.created_at DESC
			LIMIT ? OFFSET ?
		`, status, perPage, offset)
	} else {
		h.db.QueryRow("SELECT COUNT(*) FROM contact_submissions").Scan(&totalCount)
		rows, err = h.db.Query(`
			SELECT cs.id, cs.email, cs.subject, cs.message, cs.status, 
			       cs.assigned_to, u.name, cs.last_reminder_sent_at, cs.created_at, cs.updated_at
			FROM contact_submissions cs
			LEFT JOIN users u ON cs.assigned_to = u.id
			ORDER BY 
				CASE cs.status 
					WHEN 'new' THEN 1 
					WHEN 'in_progress' THEN 2 
					ELSE 3 
				END,
				cs.created_at DESC
			LIMIT ? OFFSET ?
		`, perPage, offset)
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"data":         []ContactSubmission{},
			"current_page": page,
			"last_page":    1,
			"total":        0,
		})
		return
	}
	defer rows.Close()

	submissions := []ContactSubmission{}
	for rows.Next() {
		var cs ContactSubmission
		var assignedTo sql.NullInt64
		var assignedToName, lastReminderSentAt sql.NullString

		if err := rows.Scan(
			&cs.ID, &cs.Email, &cs.Subject, &cs.Message, &cs.Status,
			&assignedTo, &assignedToName, &lastReminderSentAt, &cs.CreatedAt, &cs.UpdatedAt,
		); err != nil {
			continue
		}

		if assignedTo.Valid {
			cs.AssignedTo = &assignedTo.Int64
		}
		if assignedToName.Valid {
			cs.AssignedToName = assignedToName.String
		}
		if lastReminderSentAt.Valid {
			cs.LastReminderSentAt = &lastReminderSentAt.String
		}

		submissions = append(submissions, cs)
	}

	lastPage := (totalCount + perPage - 1) / perPage
	if lastPage < 1 {
		lastPage = 1
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":         submissions,
		"current_page": page,
		"last_page":    lastPage,
		"total":        totalCount,
	})
}

// GetContactSubmissionByID returns a single contact submission with notes
func (h *Handler) GetContactSubmissionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cs ContactSubmission
	var assignedTo sql.NullInt64
	var assignedToName, lastReminderSentAt sql.NullString

	err := h.db.QueryRow(`
		SELECT cs.id, cs.email, cs.subject, cs.message, cs.status, 
		       cs.assigned_to, u.name, cs.last_reminder_sent_at, cs.created_at, cs.updated_at
		FROM contact_submissions cs
		LEFT JOIN users u ON cs.assigned_to = u.id
		WHERE cs.id = ?
	`, id).Scan(
		&cs.ID, &cs.Email, &cs.Subject, &cs.Message, &cs.Status,
		&assignedTo, &assignedToName, &lastReminderSentAt, &cs.CreatedAt, &cs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Submission not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if assignedTo.Valid {
		cs.AssignedTo = &assignedTo.Int64
	}
	if assignedToName.Valid {
		cs.AssignedToName = assignedToName.String
	}
	if lastReminderSentAt.Valid {
		cs.LastReminderSentAt = &lastReminderSentAt.String
	}

	// Get notes
	cs.Notes = h.getContactNotes(cs.ID)

	respondJSON(w, http.StatusOK, cs)
}

// UpdateContactSubmission updates a contact submission (status, assignment)
func (h *Handler) UpdateContactSubmission(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var update struct {
		Status     string `json:"status"`
		AssignedTo *int64 `json:"assigned_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate status
	validStatuses := map[string]bool{"new": true, "in_progress": true, "resolved": true, "closed": true}
	if !validStatuses[update.Status] {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid status"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE contact_submissions SET status = ?, assigned_to = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, update.Status, update.AssignedTo, id)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update submission"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Submission updated"})
}

// AddContactNote adds an internal note to a submission
func (h *Handler) AddContactNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var note struct {
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if note.Note == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Note is required"})
		return
	}

	// Get user from context
	user, ok := getUserFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	result, err := h.db.Exec(`
		INSERT INTO contact_notes (submission_id, user_id, note)
		VALUES (?, ?, ?)
	`, id, user.UserID, note.Note)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add note"})
		return
	}

	noteID, _ := result.LastInsertId()
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":      noteID,
		"message": "Note added",
	})
}

// DeleteContactSubmission deletes a contact submission
func (h *Handler) DeleteContactSubmission(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Exec("DELETE FROM contact_submissions WHERE id = ?", id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete submission"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Submission deleted"})
}

// Helper to get contact notes
func (h *Handler) getContactNotes(submissionID int64) []ContactNote {
	rows, err := h.db.Query(`
		SELECT cn.id, cn.user_id, u.name, cn.note, cn.created_at
		FROM contact_notes cn
		JOIN users u ON cn.user_id = u.id
		WHERE cn.submission_id = ?
		ORDER BY cn.created_at DESC
	`, submissionID)

	if err != nil {
		return []ContactNote{}
	}
	defer rows.Close()

	notes := []ContactNote{}
	for rows.Next() {
		var note ContactNote
		if err := rows.Scan(&note.ID, &note.UserID, &note.UserName, &note.Note, &note.CreatedAt); err != nil {
			continue
		}
		notes = append(notes, note)
	}

	return notes
}

// Helper to get user from request context (from JWT middleware)
func getUserFromContext(r *http.Request) (middleware.UserClaims, bool) {
	return middleware.GetUserFromContext(r.Context())
}
