package controller

import (
	"Diploma/config"
	"Diploma/controller/authentication"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// 📌 Добавление экстренного контакта (Create)
func AddEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contactRequest struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&contactRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли такой пользователь в базе
	var contactUser users.User
	if err := config.DB.Where("phone = ?", contactRequest.PhoneNumber).First(&contactUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Проверяем, не добавлен ли уже этот контакт
	existingContact := users.TrustedContact{}
	if err := config.DB.Where("user_id = ? AND phone_number = ?", user.ID, contactRequest.PhoneNumber).
		First(&existingContact).Error; err == nil {
		http.Error(w, "Contact already exists", http.StatusConflict)
		return
	}

	// Добавляем экстренный контакт
	trustedContact := users.TrustedContact{
		UserID:      user.ID,
		PhoneNumber: contactRequest.PhoneNumber,
		ContactID:   contactUser.ID, // Записываем ID найденного пользователя
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&trustedContact).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s added to emergency contacts", contactRequest.PhoneNumber)
}

// 📌 Получение всех экстренных контактов пользователя (Read)
func GetEmergencyContacts(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
		http.Error(w, "Error retrieving contacts", http.StatusInternalServerError)
		return
	}

	// Структура, в которой будет и имя
	type ContactWithName struct {
		ID          uint      `json:"id"`
		UserID      uint      `json:"user_id"`
		ContactID   uint      `json:"contact_id"`
		PhoneNumber string    `json:"phone_number"`
		Name        string    `json:"name"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var result []ContactWithName

	for _, contact := range contacts {
		var contactUser users.User
		name := "Unknown"

		if err := config.DB.First(&contactUser, contact.ContactID).Error; err == nil {
			name = contactUser.Name // 🔸 используем имя из users.User
		}

		result = append(result, ContactWithName{
			ID:          contact.ID,
			UserID:      contact.UserID,
			ContactID:   contact.ContactID,
			PhoneNumber: contact.PhoneNumber,
			Name:        name,
			CreatedAt:   contact.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// 📌 Удаление экстренного контакта (Delete)
func DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contactRequest struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&contactRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var contact users.TrustedContact
	if err := config.DB.Where("user_id = ? AND phone_number = ?", user.ID, contactRequest.PhoneNumber).First(&contact).Error; err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&contact).Error; err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Emergency contact deleted successfully")
}

// 📌 Функция для проверки аутентификации пользователя
func AuthenticateUser(r *http.Request) (*users.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := authentication.ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	var user users.User
	if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}





type ContactRequest struct {
    ID        uint      `json:"id"`
    RequesterID uint    `json:"requester_id"`
    TargetID   uint     `json:"target_id"`
    Status     string   `json:"status"` // "pending", "accepted", "rejected"
    CreatedAt  time.Time `json:"created_at"`
}

// 📌 Отправка запроса на добавление контакта
func SendContactRequest(w http.ResponseWriter, r *http.Request) {
    requester, err := AuthenticateUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var request struct {
        PhoneNumber string `json:"phone_number"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Находим пользователя, которого хотят добавить
    var targetUser users.User
    if err := config.DB.Where("phone = ?", request.PhoneNumber).First(&targetUser).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Проверяем, не отправлен ли уже запрос
    var existingRequest ContactRequest
    if err := config.DB.Where("requester_id = ? AND target_id = ?", requester.ID, targetUser.ID).
        First(&existingRequest).Error; err == nil {
        http.Error(w, "Request already sent", http.StatusConflict)
        return
    }

    // Создаем запрос
    contactRequest := ContactRequest{
        RequesterID: requester.ID,
        TargetID:    targetUser.ID,
        Status:      "pending",
        CreatedAt:   time.Now(),
    }

    if err := config.DB.Create(&contactRequest).Error; err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Создаем уведомление для целевого пользователя
    notification := users.Notification{
        UserID:    targetUser.ID,
        Type:      "contact_request",
        Title:     "New Contact Request",
        Message:   fmt.Sprintf("%s wants to add you as emergency contact", requester.Name),
        CreatedAt: time.Now(),
        Metadata:  fmt.Sprintf(`{"request_id": %d}`, contactRequest.ID),
    }

    if err := config.DB.Create(&notification).Error; err != nil {
        http.Error(w, "Failed to create notification", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(contactRequest)
}

// 📌 Обработка ответа на запрос контакта
func HandleContactRequest(w http.ResponseWriter, r *http.Request) {
    user, err := AuthenticateUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var request struct {
        RequestID    uint   `json:"request_id"`
        Action       string `json:"action"` // "accept" or "reject"
        NotificationID uint `json:"notification_id"` // 👈 Добавляем ID уведомления
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Удаляем уведомление независимо от результата
    defer func() {
        if err := config.DB.Where("id = ?", request.NotificationID).Delete(&users.Notification{}).Error; err != nil {
            log.Printf("Failed to delete notification: %v", err)
        }
    }()

    // Находим запрос
    var contactRequest ContactRequest
    if err := config.DB.First(&contactRequest, request.RequestID).Error; err != nil {
        http.Error(w, "Request not found", http.StatusNotFound)
        return
    }

    // Проверяем, что пользователь является целевым
    if contactRequest.TargetID != user.ID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Обновляем статус запроса
    if request.Action == "accept" {
        // Добавляем контакты обоим пользователям
        requesterContact := users.TrustedContact{
            UserID:      contactRequest.RequesterID,
            ContactID:   contactRequest.TargetID,
            PhoneNumber: user.Phone,
            CreatedAt:   time.Now(),
        }

        targetContact := users.TrustedContact{
            UserID:      contactRequest.TargetID,
            ContactID:   contactRequest.RequesterID,
            PhoneNumber: "", // Получим из базы
            CreatedAt:   time.Now(),
        }

        // Получаем телефон реквестора
        var requester users.User
        if err := config.DB.First(&requester, contactRequest.RequesterID).Error; err != nil {
            http.Error(w, "Error finding user", http.StatusInternalServerError)
            return
        }
        targetContact.PhoneNumber = requester.Phone

        // Сохраняем контакты
        if err := config.DB.Create(&requesterContact).Error; err != nil {
            http.Error(w, "Error saving contact", http.StatusInternalServerError)
            return
        }
        if err := config.DB.Create(&targetContact).Error; err != nil {
            http.Error(w, "Error saving contact", http.StatusInternalServerError)
            return
        }

        contactRequest.Status = "accepted"
    } else {
        contactRequest.Status = "rejected"
    }

    if err := config.DB.Save(&contactRequest).Error; err != nil {
        http.Error(w, "Error updating request", http.StatusInternalServerError)
        return
    }

    // Создаем уведомление для реквестора
    var notificationMessage string
    if request.Action == "accept" {
        notificationMessage = fmt.Sprintf("%s accepted your contact request", user.Name)
    } else {
        notificationMessage = fmt.Sprintf("%s rejected your contact request", user.Name)
    }

    notification := users.Notification{
        UserID:    contactRequest.RequesterID,
        Type:      "contact_response",
        Title:     "Contact Request Update",
        Message:   notificationMessage,
        CreatedAt: time.Now(),
    }

    if err := config.DB.Create(&notification).Error; err != nil {
        http.Error(w, "Failed to create notification", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(contactRequest)
}
