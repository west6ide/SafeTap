package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
	"time"
)

func ScheduleFakeCall(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var call users.FakeCall
	if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	call.UserID = user.ID
	call.CreatedAt = time.Now()

	if err := config.DB.Create(&call).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "scheduled"})
}



func GetUserFakeCalls(w http.ResponseWriter, r *http.Request) {
	user, err := AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var calls []users.FakeCall
	if err := config.DB.Where("user_id = ?", user.ID).Order("call_time asc").Find(&calls).Error; err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(calls)
}



func DeleteFakeCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	if err := config.DB.Delete(&users.FakeCall{}, id).Error; err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
