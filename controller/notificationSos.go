package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Не указан user_id", http.StatusBadRequest)
		return
	}

	var notifications []users.Notification
	config.DB.Where("contact_id = ?", userID).Find(&notifications)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}