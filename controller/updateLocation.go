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

    // Удалим старые записи
    config.DB.Where("user_id = ?", user.ID).Delete(&users.LiveLocation{})

    // Получаем всех контактов, кому должен быть виден этот пользователь
    var contacts []users.TrustedContact
    config.DB.Where("user_id = ?", user.ID).Find(&contacts)

    // Для каждого emergency-контакта создаём видимую локацию
    for _, contact := range contacts {
        loc := users.LiveLocation{
            UserID:    user.ID,
            ContactID: contact.ContactID,
            Latitude:  req.Latitude,
            Longitude: req.Longitude,
            UpdatedAt: time.Now(),
        }
        config.DB.Create(&loc)
    }

    // Добавляем обратную видимость: пусть пользователь тоже видит контактов
    var reverseContacts []users.TrustedContact
    config.DB.Where("contact_id = ?", user.ID).Find(&reverseContacts)

    for _, rc := range reverseContacts {
        loc := users.LiveLocation{
            UserID:    user.ID,
            ContactID: rc.UserID, // обратный доступ
            Latitude:  req.Latitude,
            Longitude: req.Longitude,
            UpdatedAt: time.Now(),
        }
        config.DB.Create(&loc)
    }

    w.WriteHeader(http.StatusOK)
}



func GetEmergencyLocations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	// Получаем список contact_id для текущего пользователя
	var contactIDs []int
	if err := config.DB.Model(&users.TrustedContact{}).
		Where("user_id = ?", userID).
		Pluck("contact_id", &contactIDs).Error; err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
	}

	if len(contactIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]")) // Возвращаем пустой список
		return
	}

	// Получаем последние координаты для каждого contact_id
	var locations []users.LiveLocation
	if err := config.DB.Raw(`
		SELECT DISTINCT ON (user_id) *
		FROM live_locations
		WHERE user_id IN ?
		ORDER BY user_id, updated_at DESC
	`, contactIDs).Scan(&locations).Error; err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}



