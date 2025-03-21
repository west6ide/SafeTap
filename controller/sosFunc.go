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

    if input.Latitude == 0.0 && input.Longitude == 0.0 {
        http.Error(w, "Ошибка: Координаты не переданы или равны нулю", http.StatusBadRequest)
        return
    }

    // Логируем полученные данные
    fmt.Printf("SOS Signal received from UserID: %d, Latitude: %f, Longitude: %f\n", input.UserID, input.Latitude, input.Longitude)

    // Сохранение SOS-сигнала
    sosSignal := users.SOSSignal{
        UserID:    input.UserID,
        Latitude:  input.Latitude,
        Longitude: input.Longitude,
        CreatedAt: time.Now(),
    }
    if result := config.DB.Create(&sosSignal); result.Error != nil {
        fmt.Println("Ошибка при сохранении SOS-сигнала:", result.Error)
        http.Error(w, "Ошибка при сохранении SOS-сигнала", http.StatusInternalServerError)
        return
    }

    // Найдём контакты пользователя
    var contacts []users.TrustedContact
    config.DB.Where("UserID = ?", input.UserID).Find(&contacts)

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
