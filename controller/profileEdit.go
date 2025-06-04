package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type ProfileUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, err := authUser(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var profileUpdate ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	updateData := make(map[string]interface{})
	if profileUpdate.Name != "" {
		updateData["name"] = profileUpdate.Name
	}
	if profileUpdate.Phone != "" {
		updateData["phone"] = profileUpdate.Phone
	}
	if profileUpdate.Email != "" {
		updateData["email"] = profileUpdate.Email
	}
	if profileUpdate.Avatar != "" {
		updateData["avatar"] = profileUpdate.Avatar
	}

	if len(updateData) > 0 {
		if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
			log.Println("Failed to update profile:", err)
			http.Error(w, `{"error": "Failed to update profile"}`, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})

}
func authUser(r *http.Request) (*users.User, error) {
	var user users.User

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Расшифровка токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}
	userID := uint(userIDFloat)

	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}
