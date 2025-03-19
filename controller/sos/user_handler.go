package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"Diploma/config"
	"Diploma/users"
)

// Получаем UserID по токену
func GetUserIdHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	var user users.User
	err := config.DB.Where("token = ?", token).First(&user).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]uint{"user_id": user.ID})
}
// Получаем токен пользователя (например, по user_id)
func GetUserTokenHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var user users.User
	err := config.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": user.AccessToken})
}