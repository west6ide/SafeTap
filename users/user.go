package users

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`
	gorm.Model
	Name         string
	Phone        string
	Email        string
	Password     string `json:"-" gorm:"not null"`
	Avatar       string
	Contacts     []TrustedContact `gorm:"foreignKey:user_id"`
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	Provider     string           `json:"provider"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type TrustedContact struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	UserID      uint `gorm:"not null" json:"user_id"`
	ContactID   uint `gorm:"not null" json:"contact_id"`
	PhoneNumber string    `gorm:"not null;uniqueIndex:idx_user_phone" json:"phone_number"`
	PushToken   string    `json:"push_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}


type LiveLocation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}





type SOSSignal struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index" json:"user_id"`
    ContactID uint      `gorm:"index" json:"contact_id"`
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
    CreatedAt time.Time `json:"created_at"`
}



type Notification struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `json:"userId"`
    ContactID uint      `json:"contactId"`
    Message   string    `json:"message"`
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
	DestLatitude float64 `json:"destLatitude"`
	DestLongitude float64 `json:"destLongitude"`
    CreatedAt time.Time `json:"createdAt"`
    Type      string    `json:"type"` // "sos" или "route"
    RouteID   uint      `json:"routeId" gorm:"default:0"` // ID связанного маршрута
}



