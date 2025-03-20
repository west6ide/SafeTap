package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Обработчик регистрации Push Token
func RegisterPushTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    uint   `json:"user_id"`
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
		UserID    uint    `json:"user_id"`
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Сохранение SOS-сигнала
	sosSignal := users.SOSSignal{
		UserID:    input.UserID,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		CreatedAt: time.Now(),
	}
	config.DB.Create(&sosSignal)

	// Найдём контакты пользователя
	var contacts []users.TrustedContact
	config.DB.Where("user_id = ?", input.UserID).Find(&contacts)

	// Создаём уведомления для контактов
	for _, contact := range contacts {
		notification := users.Notification{
			UserID:    input.UserID,
			ContactID: contact.ContactID,
			Message:   "ВНИМАНИЕ! Ваш контакт отправил SOS-сигнал!",
			CreatedAt: time.Now(),
		}
		config.DB.Create(&notification)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"message": "SOS-сигнал отправлен успешно"}`)
}