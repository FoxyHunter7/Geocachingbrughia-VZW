package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	NeedsPasswordUpdate bool      `json:"needs_password_update"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUsers returns all admin users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, name, email, COALESCE(needs_password_update, 0), created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to fetch users",
		})
		return
	}
	defer rows.Close()

	users := []AdminUser{}
	for rows.Next() {
		var u AdminUser
		var needsUpdate int
		var createdAt, updatedAt string

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &needsUpdate, &createdAt, &updatedAt); err != nil {
			continue
		}

		u.NeedsPasswordUpdate = needsUpdate == 1
		u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		u.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		users = append(users, u)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   users,
		"total":  len(users),
	})
}

// CreateUser creates a new admin user and sends them an invitation email
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Invalid request body",
		})
		return
	}

	// Validate input
	errors := make(map[string][]string)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" {
		errors["name"] = append(errors["name"], "Name is required")
	}
	if req.Email == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else if !strings.Contains(req.Email, "@") {
		errors["email"] = append(errors["email"], "Invalid email format")
	}

	if len(errors) > 0 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"errors": errors,
		})
		return
	}

	// Check if email already exists
	var existingID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err == nil {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"status":  false,
			"message": "A user with this email already exists",
		})
		return
	}

	// Generate random password
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to generate password",
		})
		return
	}
	tempPassword := hex.EncodeToString(randomBytes)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to hash password",
		})
		return
	}

	// Insert user with needs_password_update = 1
	result, err := h.db.Exec(`
		INSERT INTO users (name, email, password_hash, needs_password_update, created_at, updated_at)
		VALUES (?, ?, ?, 1, datetime('now'), datetime('now'))
	`, req.Name, req.Email, string(hashedPassword))
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to create user",
		})
		return
	}

	userID, _ := result.LastInsertId()

	// Send invitation email
	h.emailService.SendAdminInvitation(req.Email, req.Name, tempPassword)

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  true,
		"message": "User created and invitation email sent",
		"data": AdminUser{
			ID:                  userID,
			Name:                req.Name,
			Email:               req.Email,
			NeedsPasswordUpdate: true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	})
}

// DeleteUser removes an admin user (cannot delete yourself)
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get current user from context to prevent self-deletion
	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if ok && currentUser.UserID == userID {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Cannot delete your own account",
		})
		return
	}

	// Check if user exists
	var existingID int64
	err = h.db.QueryRow("SELECT id FROM users WHERE id = ?", userID).Scan(&existingID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "User not found",
		})
		return
	}

	// Delete the user
	_, err = h.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to delete user",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "User deleted successfully",
	})
}

// ResendInvitation regenerates password and sends a new invitation email
func (h *Handler) ResendInvitation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get user details
	var name, email string
	var needsPasswordUpdate int
	err = h.db.QueryRow(
		"SELECT name, email, COALESCE(needs_password_update, 0) FROM users WHERE id = ?",
		userID,
	).Scan(&name, &email, &needsPasswordUpdate)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "User not found",
		})
		return
	}

	// Only allow resending if password update is still needed
	if needsPasswordUpdate != 1 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "User has already changed their password",
		})
		return
	}

	// Generate new random password
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to generate password",
		})
		return
	}
	tempPassword := hex.EncodeToString(randomBytes)

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to hash password",
		})
		return
	}

	// Update password
	_, err = h.db.Exec(
		"UPDATE users SET password = ?, updated_at = datetime('now') WHERE id = ?",
		string(hashedPassword), userID,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to update password",
		})
		return
	}

	// Send new invitation email
	h.emailService.SendAdminInvitation(email, name, tempPassword)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "New invitation email sent",
	})
}
