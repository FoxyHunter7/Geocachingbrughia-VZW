package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type AuthResponse struct {
	Status              bool   `json:"status"`
	Token               string `json:"token,omitempty"`
	User                *User  `json:"user,omitempty"`
	NeedsPasswordUpdate bool   `json:"needs_password_update,omitempty"`
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Invalid request body",
		})
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"errors": map[string][]string{
				"email":    {"Email is required"},
				"password": {"Password is required"},
			},
		})
		return
	}

	// Get user from database (including needs_password_update flag)
	var user User
	var passwordHash string
	var needsPasswordUpdate int
	err := h.db.QueryRow(
		"SELECT id, name, email, password_hash, COALESCE(needs_password_update, 0) FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &needsPasswordUpdate)

	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Invalid email or password",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Invalid email or password",
		})
		return
	}

	// Generate JWT (include needs_password_update in claims)
	token, err := h.generateJWT(user, needsPasswordUpdate == 1)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to generate token",
		})
		return
	}

	respondJSON(w, http.StatusOK, AuthResponse{
		Status:              true,
		Token:               token,
		User:                &user,
		NeedsPasswordUpdate: needsPasswordUpdate == 1,
	})
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Unauthorized",
		})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Invalid request body",
		})
		return
	}

	// Validate new password
	if len(req.NewPassword) < 8 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"errors": map[string][]string{
				"new_password": {"Password must be at least 8 characters"},
			},
		})
		return
	}

	// Get current password hash from database
	var currentHash string
	err := h.db.QueryRow(
		"SELECT password_hash FROM users WHERE id = ?",
		userClaims.UserID,
	).Scan(&currentHash)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to verify current password",
		})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Current password is incorrect",
		})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to hash new password",
		})
		return
	}

	// Update password and clear needs_password_update flag
	_, err = h.db.Exec(
		"UPDATE users SET password_hash = ?, needs_password_update = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		string(newHash), userClaims.UserID,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to update password",
		})
		return
	}

	// Generate new JWT without needs_password_update flag
	user := User{
		ID:    userClaims.UserID,
		Name:  userClaims.Name,
		Email: userClaims.Email,
	}
	token, _ := h.generateJWT(user, false)

	respondJSON(w, http.StatusOK, AuthResponse{
		Status:              true,
		Token:               token,
		User:                &user,
		NeedsPasswordUpdate: false,
	})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Unauthorized",
		})
		return
	}

	// Check if user still needs password update (from database, not just JWT)
	var needsPasswordUpdate int
	h.db.QueryRow(
		"SELECT COALESCE(needs_password_update, 0) FROM users WHERE id = ?",
		userClaims.UserID,
	).Scan(&needsPasswordUpdate)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": true,
		"data": map[string]interface{}{
			"id":                    userClaims.UserID,
			"name":                  userClaims.Name,
			"email":                 userClaims.Email,
			"needs_password_update": needsPasswordUpdate == 1,
		},
	})
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Unauthorized",
		})
		return
	}

	// Check current needs_password_update status from database
	var needsPasswordUpdate int
	h.db.QueryRow(
		"SELECT COALESCE(needs_password_update, 0) FROM users WHERE id = ?",
		userClaims.UserID,
	).Scan(&needsPasswordUpdate)

	user := User{
		ID:    userClaims.UserID,
		Name:  userClaims.Name,
		Email: userClaims.Email,
	}

	token, err := h.generateJWT(user, needsPasswordUpdate == 1)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Failed to refresh token",
		})
		return
	}

	respondJSON(w, http.StatusOK, AuthResponse{
		Status:              true,
		Token:               token,
		NeedsPasswordUpdate: needsPasswordUpdate == 1,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// With JWT, logout is handled client-side by removing the token
	// We just return success
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Logged out successfully",
	})
}

func (h *Handler) generateJWT(user User, needsPasswordUpdate bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":               user.ID,
		"email":                 user.Email,
		"name":                  user.Name,
		"needs_password_update": needsPasswordUpdate,
		"exp":                   time.Now().Add(time.Duration(h.cfg.JWT.ExpiryHours) * time.Hour).Unix(),
		"iat":                   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWT.Secret))
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
