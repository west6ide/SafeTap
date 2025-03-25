package controller

import (
	"Diploma/config"
	"Diploma/controller/authentication"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// 📌 Save SOS Signal
func SaveSOS(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUsers(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sosRequest struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sosRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve emergency contacts
	var trustedContacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&trustedContacts).Error; err != nil {
		http.Error(w, "Failed to retrieve trusted contacts", http.StatusInternalServerError)
		return
	}

	// Save each SOS signal for each trusted contact
	for _, contact := range trustedContacts {
		sosSignal := users.SOSSignal{
			UserID:    user.ID,
			ContactID: contact.ContactID,
			Latitude:  sosRequest.Latitude,
			Longitude: sosRequest.Longitude,
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&sosSignal).Error; err != nil {
			http.Error(w, "Failed to save SOS signal", http.StatusInternalServerError)
			return
		}

		// Save notification
		notification := users.Notification{
			UserID:    contact.ContactID,
			ContactID: user.ID,
			Message:   fmt.Sprintf("SOS signal received from %s", user.Name),
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&notification).Error; err != nil {
			http.Error(w, "Failed to save notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"SOS signal saved successfully"}`)
}

func authenticateUsers(r *http.Request) (*users.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("Missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := authentication.ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	var user users.User
	if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("User not found")
	}

	return &user, nil
}
