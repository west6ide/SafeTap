// controller/location.go
package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Обновление координат текущего пользователя
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

	location := users.LiveLocation{
		UserID:    user.ID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		UpdatedAt: time.Now(),
	}

	config.DB.
		Where("user_id = ?", user.ID).
		Assign(location).
		FirstOrCreate(&location)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Location updated")
}


// Получение координат emergency-контактов
func GetEmergencyLocations(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
		http.Error(w, "Failed to get contacts", http.StatusInternalServerError)
		return
	}

	var contactIDs []uint
	for _, c := range contacts {
		contactIDs = append(contactIDs, c.ContactID)
	}
	contactIDs = append(contactIDs, user.ID) // Чтобы пользователь тоже видел себя

	var locations []users.LiveLocation
	if err := config.DB.
		Where("user_id IN ?", contactIDs).
		Find(&locations).Error; err != nil {
		http.Error(w, "Failed to get locations", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(locations)
}

