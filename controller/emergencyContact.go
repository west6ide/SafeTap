package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ContactRequest struct {
	PhoneNumber string `json:"phone_number"`
	PushToken   string `json:"push_token"`
}

func AddEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contactRequest ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&contactRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	existingContact := users.TrustedContact{}
	if err := config.DB.Where("user_id = ? AND phone_number = ?", user.ID, contactRequest.PhoneNumber).First(&existingContact).Error; err == nil {
		http.Error(w, "Contact already exists", http.StatusConflict)
		return
	}

	trustedContact := users.TrustedContact{
		UserID:      user.ID,
		PhoneNumber: contactRequest.PhoneNumber,
		PushToken:   contactRequest.PushToken,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&trustedContact).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Emergency contact added successfully")
}

func DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contactRequest ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&contactRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := config.DB.Where("user_id = ? AND phone_number = ?", user.ID, contactRequest.PhoneNumber).Delete(&users.TrustedContact{}).Error; err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Emergency contact deleted successfully")
}

func authenticateUser(r *http.Request) (*users.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	var user users.User
	if err := config.DB.Where("access_token = ?", token).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	return &user, nil
}
