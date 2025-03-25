package controller

import (
	"Diploma/config"
	"Diploma/users"
	"Diploma/controller/authentication"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// 📌 Функция для аутентификации пользователя
func authenticateUsers(r *http.Request) (*users.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := authentication.ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	var user users.User
	if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// 📌 Endpoint for receiving and saving SOS signals
func SendSOS(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUsers(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Fetch emergency contacts of the user
	var contacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
		http.Error(w, `{"error": "Failed to retrieve contacts"}`, http.StatusInternalServerError)
		return
	}

	// Save SOS signals for each contact and create notifications
	for _, contact := range contacts {
		// Save the SOS signal
		sosSignal := users.SOSSignal{
			UserID:    user.ID,
			ContactID: contact.ContactID,
			Latitude:  requestBody.Latitude,
			Longitude: requestBody.Longitude,
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&sosSignal).Error; err != nil {
			http.Error(w, `{"error": "Failed to save SOS signal"}`, http.StatusInternalServerError)
			return
		}

		// Create notification for the contact
		notification := users.Notification{
			UserID:    contact.ContactID,
			ContactID: user.ID,
			Message:   "SOS signal received!",
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&notification).Error; err != nil {
			http.Error(w, `{"error": "Failed to save notification"}`, http.StatusInternalServerError)
			return
		}
	}

	// Send JSON response on success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "SOS signals sent successfully!"}`))
}
