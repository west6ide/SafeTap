package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
	"fmt"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
    user, err := authenticateUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        http.Error(w, "User ID не передан", http.StatusBadRequest)
        return
    }

    // Проверка, что текущий пользователь имеет право получать уведомления этого user_id
    if fmt.Sprint(user.ID) != userID {
        http.Error(w, "Access denied: You are not allowed to view these notifications", http.StatusForbidden)
        return
    }

    var notifications []users.Notification
    result := config.DB.Where("user_id = ?", userID).Find(&notifications)
    if result.Error != nil {
        http.Error(w, "Ошибка получения уведомлений", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(notifications)
}





