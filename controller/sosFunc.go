package controller

import (
	"Diploma/config"
	"Diploma/users"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// Обработчик отправки SOS-сигнала
func SendSOSHandler(w http.ResponseWriter, r *http.Request) {
	var sos users.SOSSignal
	json.NewDecoder(r.Body).Decode(&sos)
	sos.CreatedAt = time.Now()
	config.DB.Create(&sos)

	var contacts []users.TrustedContact
	config.DB.Where("user_id = ?", sos.UserID).Find(&contacts)

	for _, contact := range contacts {
		if contact.PushToken != "" {
			message := fmt.Sprintf("SOS! Ваш контакт в опасности. Координаты: %.6f, %.6f", sos.Latitude, sos.Longitude)
			go sendPushNotification(contact.PushToken, message)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Функция отправки push-уведомлений через Firebase
func sendPushNotification(token, message string) {
	fcmKey := os.Getenv("FCM_SERVER_KEY")
	payload := map[string]interface{}{
		"to": token,
		"notification": map[string]string{
			"title": "SOS Сигнал!",
			"body": message,
		},
	}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "key="+fcmKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Ответ Firebase:", string(body))
}