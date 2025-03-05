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
	log.Println("🔄 Новое WebSocket-подключение")
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

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket ошибка:", err)
		return
	}
	defer ws.Close()

	log.Printf("✅ Подключение установлено для пользователя ID=%d\n", claims.UserID)

	// Добавляем WebSocket клиента
	locationMu.Lock()
	clients[ws] = claims.UserID
	locationMu.Unlock()

	// Ping-Pong для поддержания соединения
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Обработка сообщений WebSocket
	for {
		var loc LocationUpdate
		err := ws.ReadJSON(&loc)
		if err != nil {
			log.Println("❌ Ошибка чтения WebSocket:", err)
			break
		}

		// Сохранение координат в БД
		location := users.LiveLocation{
			UserID:    claims.UserID,
			Lat:       loc.Lat,
			Lng:       loc.Lng,
			UpdatedAt: time.Now(),
		}
		if err := config.DB.Save(&location).Error; err != nil {
			log.Println("❌ Ошибка сохранения локации:", err)
			continue
		}

		log.Printf("✅ Локация обновлена: ID=%d, Lat=%.6f, Lng=%.6f\n", claims.UserID, loc.Lat, loc.Lng)

		// Рассылка локации экстренным контактам
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

	log.Printf("🔄 Найдено %d координат для отправки WebSocket клиентам", len(locations))

	// Отправляем данные всем подключённым клиентам
	locationJSON, _ := json.Marshal(locations)
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, locationJSON)
		if err != nil {
			log.Println("❌ WebSocket ошибка отправки:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
