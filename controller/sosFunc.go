package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SOSRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func SendSOS(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sosRequest struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sosRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch emergency contacts
	var trustedContacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&trustedContacts).Error; err != nil {
		http.Error(w, "Error fetching emergency contacts", http.StatusInternalServerError)
		return
	}

	// Save SOS signals and create notifications for each contact
	for _, contact := range trustedContacts {
		// Save SOS signal
		sosSignal := users.SOSSignal{
			UserID:    user.ID,
			ContactID: contact.ContactID,
			Latitude:  sosRequest.Latitude,
			Longitude: sosRequest.Longitude,
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&sosSignal).Error; err != nil {
			http.Error(w, "Error saving SOS signal", http.StatusInternalServerError)
			return
		}

		// Save Notification
		notification := users.Notification{
			UserID:    contact.ContactID,
			ContactID: user.ID,
			Message:   "SOS signal received!",
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&notification).Error; err != nil {
			http.Error(w, "Error saving notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "SOS signals sent successfully!")
}
