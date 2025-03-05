package controller

import (
	"Diploma/config"
	"Diploma/controller/authentication"
	"Diploma/users"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var clients = make(map[*websocket.Conn]uint)
var locationMu sync.Mutex

type LocationUpdate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func HandleLiveLocation(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Println("❌ Ошибка: JWT не передан")
		http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
		return
	}

	claims, err := authentication.ValidateJWT(token)
	if err != nil {
		log.Println("❌ Ошибка JWT:", err)
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	// Убеждаемся, что ошибок нет перед установкой WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket ошибка:", err)
		return
	}
	defer ws.Close()

	log.Printf("✅ Подключение пользователя %d к WebSocket\n", claims.UserID)

	locationMu.Lock()
	clients[ws] = claims.UserID
	locationMu.Unlock()

	for {
		var loc LocationUpdate
		err := ws.ReadJSON(&loc)
		if err != nil {
			log.Println("❌ Ошибка чтения WebSocket:", err)
			break
		}

		config.DB.Save(&users.LiveLocation{
			UserID:    claims.UserID,
			Lat:       loc.Lat,
			Lng:       loc.Lng,
			UpdatedAt: time.Now(),
		})

		broadcastLocation(claims.UserID)
	}
}

// Передаёт координаты экстренным контактам
func broadcastLocation(userID uint) {
	locationMu.Lock()
	defer locationMu.Unlock()

	// Получаем список экстренных контактов
	var contacts []users.TrustedContact
	config.DB.Where("user_id = ?", userID).Find(&contacts)

	var contactIDs []uint
	for _, contact := range contacts {
		contactID, err := strconv.ParseUint(contact.ContactID, 10, 32)
		if err != nil {
			log.Println("Ошибка преобразования ContactID:", err)
			continue
		}
		contactIDs = append(contactIDs, uint(contactID))
	}

	// Получаем местоположение всех контактов
	var locations []users.LiveLocation
	config.DB.Where("user_id IN ?", contactIDs).Find(&locations)

	// Отправляем всем подключённым клиентам
	locationJSON, _ := json.Marshal(locations)
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, locationJSON)
		if err != nil {
			log.Println("WebSocket error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
