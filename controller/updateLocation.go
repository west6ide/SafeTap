// controller/location.go
package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
)

func UpdateLiveLocation(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var location users.LiveLocation
	err = config.DB.Where("user_id = ?", user.ID).First(&location).Error
	if err != nil {
		location = users.LiveLocation{
			UserID:   user.ID,
			Latitude: req.Latitude,
			Longitude: req.Longitude,
		}
		config.DB.Create(&location)
	} else {
		location.Latitude = req.Latitude
		location.Longitude = req.Longitude
		config.DB.Save(&location)
	}

	w.WriteHeader(http.StatusOK)
}


func GetEmergencyContactsLocations(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var contactIDs []uint
	for _, c := range contacts {
		contactIDs = append(contactIDs, c.ContactID)
	}

	var locations []users.LiveLocation
	if len(contactIDs) > 0 {
		if err := config.DB.Where("user_id IN ?", contactIDs).Find(&locations).Error; err != nil {
			http.Error(w, "Failed to retrieve locations", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
