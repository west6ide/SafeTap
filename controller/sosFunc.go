package controller

import (
	"Diploma/config"
	"Diploma/controller/authentication"
	"Diploma/users"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]uint) // Хранение клиентов по UserID
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

	// Поддержка ping/pong для предотвращения разрыва соединения
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Чтение сообщений (координат) от клиента
	for {
		var loc LocationUpdate
		err := ws.ReadJSON(&loc)
		if err != nil {
			log.Println("❌ Ошибка чтения WebSocket:", err)
			break
		}

		var location users.LiveLocation

		// Проверяем, есть ли уже координаты пользователя
		result := config.DB.Where("user_id = ?", claims.UserID).First(&location)
		if result.RowsAffected == 0 {
			// Если записи нет, создаём новую
			location = users.LiveLocation{
				UserID:    claims.UserID,
				Lat:       loc.Lat,
				Lng:       loc.Lng,
				UpdatedAt: time.Now(),
			}
			if err := config.DB.Create(&location).Error; err != nil {
				log.Println("❌ Ошибка создания локации:", err)
				continue
			}
			log.Printf("✅ Локация создана: ID=%d, Lat=%.6f, Lng=%.6f\n", claims.UserID, loc.Lat, loc.Lng)
		} else {
			// Если запись есть, обновляем координаты
			location.Lat = loc.Lat
			location.Lng = loc.Lng
			location.UpdatedAt = time.Now()

			if err := config.DB.Save(&location).Error; err != nil {
				log.Println("❌ Ошибка обновления локации:", err)
				continue
			}
			log.Printf("✅ Локация обновлена: ID=%d, Lat=%.6f, Lng=%.6f\n", claims.UserID, loc.Lat, loc.Lng)
		}

		// Рассылка обновленных координат
		broadcastLocation(claims.UserID)
	}
}

// Рассылает данные всем экстренным контактам
func broadcastLocation(userID uint) {
	locationMu.Lock()
	defer locationMu.Unlock()

	var contacts []users.TrustedContact
	config.DB.Where("user_id = ?", userID).Find(&contacts)

	var location users.LiveLocation
	config.DB.Where("user_id = ?", userID).First(&location)

	locationJSON, _ := json.Marshal(location)
	log.Printf("📡 Отправка координат %d клиентам\n", len(clients))

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, locationJSON)
		if err != nil {
			log.Println("❌ WebSocket ошибка отправки:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
