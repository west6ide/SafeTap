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
	ContactID   string
	PhoneNumber string    `gorm:"not null;uniqueIndex:idx_user_phone" json:"phone_number"`
	PushToken   string    `json:"push_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type LiveLocation struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"` // ID пользователя
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	UpdatedAt time.Time `json:"updated_at"`
}
