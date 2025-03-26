package controller

import (
	"Diploma/config"
	"Diploma/controller/authentication"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func SaveSOS(w http.ResponseWriter, r *http.Request) {
    user, err := authenticateUsers(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var requestBody struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    }

    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    var contacts []users.TrustedContact
    if err := config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
        http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
        return
    }

    if len(contacts) == 0 {
        http.Error(w, "No trusted contacts found", http.StatusNotFound)
        return
    }

    fmt.Println("✅ Retrieved Contacts: ", contacts)

    for _, contact := range contacts {
        sosSignal := users.SOSSignal{
            UserID:    user.ID,
            ContactID: contact.ContactID,
            Latitude:  requestBody.Latitude,
            Longitude: requestBody.Longitude,
            CreatedAt: time.Now(),
        }

        fmt.Printf("📌 Saving SOS Signal for ContactID: %d\n", contact.ContactID)

        if err := config.DB.Create(&sosSignal).Error; err != nil {
            fmt.Println("❌ Error saving SOSSignal:", err)
            continue
        }

        fmt.Println("✅ Successfully saved SOSSignal for ContactID:", contact.ContactID)

        notification := users.Notification{
            UserID:    contact.ContactID,
            ContactID: user.ID,
            Message:   fmt.Sprintf("SOS signal received from UserID %d", user.ID),
            CreatedAt: time.Now(),
        }

        fmt.Printf("📌 Saving Notification for ContactID: %d\n", contact.ContactID)

        if err := config.DB.Create(&notification).Error; err != nil {
            fmt.Println("❌ Error saving Notification:", err)
            continue
        }

        fmt.Println("✅ Successfully saved Notification for ContactID:", contact.ContactID)
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "SOS signal processed successfully")
}


func authenticateUsers(r *http.Request) (*users.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("Missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := authentication.ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	var user users.User
	if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("User not found")
	}

	return &user, nil
}
