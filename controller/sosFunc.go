package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
	"time"
)

// Обработчик регистрации Push Token
func RegisterPushTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    uint   `json:"user_id"`
		PushToken string `json:"push_token"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	config.DB.Model(&users.User{}).Where("id = ?", request.UserID).Update("push_token", request.PushToken)
	w.WriteHeader(http.StatusOK)
}

func SendSOS(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sosRequest struct {
		UserID    uint    `json:"user_id"`
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sosRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sosSignal := users.SOSSignal{
		UserID:    user.ID,
		Latitude:  sosRequest.Latitude,
		Longitude: sosRequest.Longitude,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&sosSignal).Error; err != nil {
		http.Error(w, "Failed to save SOS signal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("SOS signal saved successfully"))
}

