// controller/location.go
package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
	"time"
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


// func GetEmergencyContactsLocations(w http.ResponseWriter, r *http.Request) {
// 	user, err := AuthenticateUser(r)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var contacts []users.TrustedContact
// 	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}

// 	var contactIDs []uint
// 	for _, c := range contacts {
// 		contactIDs = append(contactIDs, c.ContactID)
// 	}

// 	var locations []users.LiveLocation
// 	if len(contactIDs) > 0 {
// 		if err := config.DB.Where("user_id IN ?", contactIDs).Find(&locations).Error; err != nil {
// 			http.Error(w, "Failed to retrieve locations", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(locations)
// }

func GetEmergencyLocations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var contactIDs []int
	if err := config.DB.Model(&users.TrustedContact{}).
		Where("user_id = ?", userID).
		Pluck("contact_id", &contactIDs).Error; err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
	}

	if len(contactIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	type LocationWithName struct {
		ID        uint      `json:"id"`
		UserID    int       `json:"user_id"`
		ContactID int       `json:"contact_id"`
		Latitude  float64   `json:"latitude"`
		Longitude float64   `json:"longitude"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"` // 👈 новое поле
	}

	var results []LocationWithName
	if err := config.DB.Raw(`
		SELECT DISTINCT ON (l.user_id) l.*, u.name
		FROM live_locations l
		JOIN users u ON l.user_id = u.id
		WHERE l.user_id IN ?
		ORDER BY l.user_id, l.updated_at DESC
	`, contactIDs).Scan(&results).Error; err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

