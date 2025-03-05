package controller

import (
	"Diploma/config"
	"Diploma/controller/authentication"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 📌 Добавление экстренного контакта (Create)
func AddEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
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
		ContactID:   strconv.Itoa(int(contactUser.ID)), // Записываем ID найденного пользователя
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
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var contacts []users.TrustedContact
	if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
		http.Error(w, "Error retrieving contacts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

// 📌 Обновление номера экстренного контакта (Update)
func UpdateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateRequest struct {
		OldPhoneNumber string `json:"old_phone_number"`
		NewPhoneNumber string `json:"new_phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var contact users.TrustedContact
	if err := config.DB.Where("user_id = ? AND phone_number = ?", user.ID, updateRequest.OldPhoneNumber).First(&contact).Error; err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	// Проверяем, есть ли новый номер в базе пользователей
	var newContactUser users.User
	if err := config.DB.Where("phone = ?", updateRequest.NewPhoneNumber).First(&newContactUser).Error; err != nil {
		http.Error(w, "New phone number not found in users", http.StatusNotFound)
		return
	}

	// Обновляем номер телефона и ID пользователя
	contact.PhoneNumber = updateRequest.NewPhoneNumber
	contact.ContactID = strconv.Itoa(int(newContactUser.ID))

	if err := config.DB.Save(&contact).Error; err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Emergency contact updated successfully")
}

// 📌 Удаление экстренного контакта (Delete)
func DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(r)
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
func authenticateUser(r *http.Request) (*users.User, error) {
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
