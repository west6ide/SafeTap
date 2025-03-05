package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients    = make(map[*websocket.Conn]uint) // uint вместо string
	locationMu sync.Mutex
)

type LocationUpdate struct {
	UserID uint    `json:"user_id"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

func HandleLiveLocation(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	locationMu.Lock()
	clients[ws] = user.ID
	locationMu.Unlock()

	for {
		var loc LocationUpdate
		if err := ws.ReadJSON(&loc); err != nil {
			log.Println("Error reading location:", err)
			break
		}

		// Сохраняем координаты в БД
		config.DB.Save(&users.LiveLocation{
			UserID:    user.ID,
			Lat:       loc.Lat,
			Lng:       loc.Lng,
			UpdatedAt: time.Now(),
		})

		// Отправляем обновления экстренным контактам
		broadcastLocation(user.ID)
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
