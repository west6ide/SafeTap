package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Обработчик регистрации Push Token
func RegisterPushTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    uint   `json:"UserID"`
		PushToken string `json:"push_token"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	config.DB.Model(&users.User{}).Where("id = ?", request.UserID).Update("push_token", request.PushToken)
	w.WriteHeader(http.StatusOK)
}

func SendSOS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		UserID    uint    `json:"UserID"`
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	log.Printf("Получен SOS-запрос: UserID=%d, Latitude=%.6f, Longitude=%.6f", input.UserID, input.Latitude, input.Longitude)
	log.Printf("Запрос в формате JSON: %v", input)

	sosSignal := users.SOSSignal{
		UserID:    input.UserID,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		CreatedAt: time.Now(),
	}
	if err := config.DB.Create(&sosSignal).Error; err != nil {
		log.Printf("Ошибка при сохранении SOS-сигнала: %v", err)
		http.Error(w, "Ошибка сохранения SOS-сигнала", http.StatusInternalServerError)
		return
	}

	var contacts []users.TrustedContact
	config.DB.Where("UserID = ?", input.UserID).Find(&contacts)

	for _, contact := range contacts {
		notification := users.Notification{
			UserID:    input.UserID,
			ContactID: contact.ContactID,
			Message:   "ВНИМАНИЕ! Ваш контакт отправил SOS-сигнал!",
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&notification).Error; err != nil {
			log.Printf("Ошибка при сохранении уведомления: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"message": "SOS-сигнал отправлен успешно"}`)
}
