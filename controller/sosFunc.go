package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
	"time"
)

// SaveSOS handles saving an SOS signal to the database
func SaveSOS(w http.ResponseWriter, r *http.Request) {
    user, err := AuthenticateUser(r)
    if err != nil {
        http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
        return
    }

    var sosRequest struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    }

    if err := json.NewDecoder(r.Body).Decode(&sosRequest); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var contacts []users.TrustedContact
    if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    for _, contact := range contacts {
        sosSignal := users.SOSSignal{
            UserID:    user.ID,
            ContactID: contact.ContactID,
            Latitude:  sosRequest.Latitude,
            Longitude: sosRequest.Longitude,
            CreatedAt: time.Now(),
        }
        if err := config.DB.Create(&sosSignal).Error; err != nil {
            http.Error(w, "Failed to save SOS signal: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "SOS signals saved successfully"})
}
