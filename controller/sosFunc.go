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

// Отправка SOS
func SendSOS(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sosRequest SOSRequest
	if err := json.NewDecoder(r.Body).Decode(&sosRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получение всех экстренных контактов пользователя
	var contacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
		http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
		return
	}

	// Сохранение SOS сигнала для каждого контакта
	for _, contact := range contacts {
		// Сохранение SOS сигнала
		signal := users.SOSSignal{
			UserID:    user.ID,
			ContactID: contact.ContactID,
			Latitude:  sosRequest.Latitude,
			Longitude: sosRequest.Longitude,
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&signal).Error; err != nil {
			http.Error(w, "Failed to save SOS signal", http.StatusInternalServerError)
			return
		}

		// Сохранение уведомления
		notification := users.Notification{
			UserID:    contact.ContactID,
			ContactID: user.ID,
			Message:   fmt.Sprintf("SOS сигнал от пользователя %d", user.ID),
			CreatedAt: time.Now(),
		}

		if err := config.DB.Create(&notification).Error; err != nil {
			http.Error(w, "Failed to save notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SOS сигнал успешно отправлен!"))
}
