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



func GetLiveLocations(w http.ResponseWriter, r *http.Request) {
    user, err := AuthenticateUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var locations []users.LiveLocation
    config.DB.
        Where("contact_id = ?", user.ID).
        Order("updated_at desc").
        Find(&locations)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(locations)
}


