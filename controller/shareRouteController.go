package controller

import (
	"Diploma/users"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type SharedRouteRequest struct {
    SenderID  uint    `json:"senderId"`
    StartLat  float64 `json:"startLat"`
    StartLng  float64 `json:"startLng"`
    DestLat   float64 `json:"destLat"`
    DestLng   float64 `json:"destLng"`
    Duration  string  `json:"duration"`
    Distance  string  `json:"distance"`
}

// POST /share_route
func ShareRouteHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req SharedRouteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        route := users.SharedRoute{
            SenderID: req.SenderID,
            StartLat: req.StartLat,
            StartLng: req.StartLng,
            DestLat:  req.DestLat,
            DestLng:  req.DestLng,
            Duration: req.Duration,
            Distance: req.Distance,
        }

        if err := db.Create(&route).Error; err != nil {
            http.Error(w, "Failed to save route", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"status": "Route shared"})
    }
}

// GET /shared_routes?userId=123
func GetSharedRoutesHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, err := ParseUintQueryParam(r, "userId")
        if err != nil {
            http.Error(w, "Invalid userId", http.StatusBadRequest)
            return
        }

        var routes []users.SharedRoute
        if err := db.Where("sender_id = ?", userID).Find(&routes).Error; err != nil {
            http.Error(w, "Failed to fetch routes", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(routes)
    }
}
