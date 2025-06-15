package controller

import (
	"Diploma/users"
	"encoding/json"
	"fmt"
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

        // Создаем маршрут
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

        // Создаем уведомление о маршруте
        notification := users.Notification{
            UserID:    route.SenderID,
            Message:   fmt.Sprintf("Shared route: %s (%s)", route.Distance, route.Duration),
            Latitude:  route.StartLat,
            Longitude: route.StartLng,
            Type:      "route",
            RouteID:   route.ID,
        }

        if err := db.Create(&notification).Error; err != nil {
            http.Error(w, "Failed to create notification", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"status": "Route shared"})
    }
}

// GET /shared_routes?userId=123
func GetSharedRoutesHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := AuthenticateUser(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var routes []users.SharedRoute
        if err := db.Where("sender_id = ?", user.ID).Find(&routes).Error; err != nil {
            http.Error(w, "Failed to fetch routes", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(routes)
    }
}
