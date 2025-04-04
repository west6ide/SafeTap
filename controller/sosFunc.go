package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SaveSOS(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
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

	// Retrieve all emergency contacts for this user
	var emergencyContacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&emergencyContacts).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	for _, contact := range emergencyContacts {
		// Save the SOS Signal
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

		// Create Notification for each contact
		notification := users.Notification{
			UserID:    contact.ContactID,
			ContactID: user.ID,
			Message:   fmt.Sprintf("Your trusted contact %s has sent an SOS signal!", user.Name),
			Latitude:  sosRequest.Latitude,
			Longitude: sosRequest.Longitude,
			CreatedAt: time.Now(),
		}
		

		if err := config.DB.Create(&notification).Error; err != nil {
			http.Error(w, "Failed to save notification", http.StatusInternalServerError)
			return
		}
		deleteNotificationAfterDelay(notification.ID, 30*time.Minute)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "SOS signal sent successfully!")
}

func deleteNotificationAfterDelay(notificationID uint, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		config.DB.Delete(&users.Notification{}, notificationID)
	}()
}

func StartNotificationCleaner() {
	go func() {
		for {
			// Delete notifications older than 30 minutes
			expirationTime := time.Now().Add(-30 * time.Minute)
			config.DB.Where("created_at < ?", expirationTime).Delete(&users.Notification{})

			time.Sleep(5 * time.Minute) // Run cleanup every 5 minutes
		}
	}()
}

