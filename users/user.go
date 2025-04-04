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
	ID        uint      `json:"ID"`
	UserID    uint       `json:"UserID"`
	ContactID uint       `json:"ContactID"`
	Latitude  float64   `json:"Latitude"`
	Longitude float64   `json:"Longitude"`
	UpdatedAt time.Time `json:"UpdatedAt"`
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
	UserID    uint      `gorm:"index" json:"user_id"`        // User who receives the notification
	ContactID uint      `gorm:"index" json:"contact_id"`     // User who sent the SOS signal
	Message   string    `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}



