package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ProfileUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, err := authUser(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var profileUpdate ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	updateData := make(map[string]interface{})
	if profileUpdate.Name != "" {
		updateData["name"] = profileUpdate.Name
	}
	if profileUpdate.Phone != "" {
		updateData["phone"] = profileUpdate.Phone
	}
	if profileUpdate.Email != "" {
		updateData["email"] = profileUpdate.Email
	}
	if profileUpdate.Avatar != "" {
		updateData["avatar"] = profileUpdate.Avatar
	}

	if len(updateData) > 0 {
		if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
			log.Println("Failed to update profile:", err)
			http.Error(w, `{"error": "Failed to update profile"}`, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})

}
func authUser(r *http.Request) (*users.User, error) {
	var user users.User

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if err := config.DB.Where("access_token = ?", token).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invailed token")
	}
	return &user, nil

}
