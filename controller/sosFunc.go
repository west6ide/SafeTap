package controller

import (
	"Diploma/config"
	"Diploma/users"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Connection struct {
	WS       *websocket.Conn
	UserID   string
	Contacts []string
}

var (
	connections = make(map[string]*Connection) // userID -> connection
	mu          sync.Mutex
)

// Обработчик подключения
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket:", err)
		return
	}
	defer ws.Close()

	var conn *Connection

	for {
		var loc users.Location
		if err := ws.ReadJSON(&loc); err != nil {
			log.Println("Read error:", err)
			break
		}

		mu.Lock()
		if conn == nil {
			// При первом сообщении регистрируем пользователя
			contacts := getUserContacts(loc.UserID)
			conn = &Connection{
				WS:       ws,
				UserID:   loc.UserID,
				Contacts: contacts,
			}
			connections[loc.UserID] = conn
			log.Printf("User %s connected with %d contacts\n", loc.UserID, len(contacts))
		}
		mu.Unlock()

		// Сохраняем координаты (если нужно)
		saveLocation(loc)

		// Рассылаем обновление контактам
		broadcastLocationToContacts(loc)
	}

	mu.Lock()
	delete(connections, conn.UserID)
	mu.Unlock()
}

// Получаем список контактов для пользователя
func getUserContacts(userID string) []string {
	var contacts []string
	rows, err := config.DB.Raw(`
        SELECT u.username 
        FROM emergency_contacts ec
        JOIN users u ON ec.contact_id = u.id
        WHERE ec.user_id = (SELECT id FROM users WHERE username = ?)
    `, userID).Rows()
	if err != nil {
		log.Printf("Failed to fetch contacts for %s: %v\n", userID, err)
		return contacts
	}
	defer rows.Close()

	for rows.Next() {
		var contact string
		rows.Scan(&contact)
		contacts = append(contacts, contact)
	}
	return contacts
}

// Сохраняем локацию в базу (если нужно)
func saveLocation(loc users.Location) {
	config.DB.Exec("INSERT INTO locations (user_id, lat, lng, created_at) VALUES ((SELECT id FROM users WHERE username = ?), ?, ?, ?)",
		loc.UserID, loc.Lat, loc.Lng, time.Now())
}

// Оповещаем все контакты о новой локации
func broadcastLocationToContacts(loc users.Location) {
	mu.Lock()
	defer mu.Unlock()

	if conn, exists := connections[loc.UserID]; exists {
		for _, contact := range conn.Contacts {
			if contactConn, ok := connections[contact]; ok {
				contactConn.WS.WriteJSON(loc)
			}
		}
	}
}
