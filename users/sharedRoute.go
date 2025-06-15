package users

import "time"

type SharedRoute struct {
    ID        uint      `gorm:"primaryKey"`
    SenderID  uint      `json:"senderId"`
    StartLat  float64   `json:"startLat"`
    StartLng  float64   `json:"startLng"`
    DestLat   float64   `json:"destLat"`
    DestLng   float64   `json:"destLng"`
    Duration  string    `json:"duration"`
    Distance  string    `json:"distance"`
    CreatedAt time.Time `json:"createdAt"`
}
