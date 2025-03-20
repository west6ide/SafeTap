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
	Contacts     []TrustedContact `gorm:"foreignKey:UserID"`
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
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index" json:"user_id"` // ID пользователя
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	UpdatedAt time.Time `json:"updated_at"`
}
// Модель SOS-сигнала
type SOSSignal struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Latitude  float64   `gorm:"not null"`
	Longitude float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Модель уведомления
type Notification struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ContactID uint      `gorm:"not null"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
