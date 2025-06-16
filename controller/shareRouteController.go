package controller

import (
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type SharedRouteRequest struct {
    StartLat  float64 `json:"startLat"`
    StartLng  float64 `json:"startLng"`
    DestLat   float64 `json:"destLat"`
    DestLng   float64 `json:"destLng"`
    Duration  string  `json:"duration"`
    Distance  string  `json:"distance"`
}

func ShareRouteHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := AuthenticateUser(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var req SharedRouteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Create the route
        route := users.SharedRoute{
            SenderID: user.ID,
            StartLat: req.StartLat,
            StartLng: req.StartLng,
            DestLat:  req.DestLat,
            DestLng:  req.DestLng,
            Duration: req.Duration,
            Distance: req.Distance,
            CreatedAt: time.Now(),
        }

        if err := db.Create(&route).Error; err != nil {
            http.Error(w, "Failed to save route", http.StatusInternalServerError)
            return
        }

        // Get user's contacts
        var contacts []users.TrustedContact
        if err := db.Where("user_id = ?", user.ID).Find(&contacts).Error; err != nil {
            http.Error(w, "Failed to get contacts", http.StatusInternalServerError)
            return
        }

        // Create notifications for each contact
        for _, contact := range contacts {
            notification := users.Notification{
                UserID:      contact.ContactID,
                ContactID:   user.ID,
                Message:     fmt.Sprintf("%s shared a route with you: %s (%s)", user.Name, route.Distance, route.Duration),
                Latitude:    route.StartLat,    // Start point latitude
                Longitude:   route.StartLng,    // Start point longitude
                DestLatitude:  route.DestLat,   // Destination latitude
                DestLongitude: route.DestLng,   // Destination longitude
                Type:        "route",
                RouteID:     route.ID,
                CreatedAt:   time.Now(),
            }

            if err := db.Create(&notification).Error; err != nil {
                http.Error(w, "Failed to create notification", http.StatusInternalServerError)
                return
            }
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"status": "Route shared with all contacts"})
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
