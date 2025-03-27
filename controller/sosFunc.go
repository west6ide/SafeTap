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

	var request struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var trustedContacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&trustedContacts).Error; err != nil {
		http.Error(w, "Error fetching contacts", http.StatusInternalServerError)
		return
	}

	for _, contact := range trustedContacts {
		// Save the SOS signal in the database
		sosSignal := users.SOSSignal{
			UserID:    user.ID,
			ContactID: contact.ContactID,
			Latitude:  request.Latitude,
			Longitude: request.Longitude,
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&sosSignal).Error; err != nil {
			http.Error(w, "Error saving SOS signal", http.StatusInternalServerError)
			return
		}

		// Save the notification for each contact
		notification := users.Notification{
			UserID:    contact.ContactID,
			ContactID: user.ID,
			Message:   "SOS Signal received!",
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&notification).Error; err != nil {
			http.Error(w, "Error saving notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "SOS signals and notifications successfully saved.")
}
