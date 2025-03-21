package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        http.Error(w, "User ID не передан", http.StatusBadRequest)
        return
    }

    var notifications []users.Notification
    result := config.DB.Where("contact_id = ?", userID).Find(&notifications)
    if result.Error != nil {
        http.Error(w, "Ошибка получения уведомлений", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(notifications)
}
