package controllers

import (
	"backend-test-mekari/internal/services"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	authSvc services.AuthService
}

func NewAuthController(as services.AuthService) *AuthController {
	return &AuthController{authSvc: as}
}

// HandleLogin godoc
// @Summary User Login
// @Description Login untuk mendapatkan token. Gunakan email: employee@mekari.com atau manager@mekari.com
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Credentials"
// @Success 200 {object} models.LoginResponse
// @Router /api/auth/login [post]
func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, role, userID, err := c.authSvc.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   token,
		"role":    role,
		"user_id": userID,
	})
}
