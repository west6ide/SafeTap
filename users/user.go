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
	Password     string           `json:"-" gorm:"not null"`
	Contacts     []TrustedContact `gorm:"foreignKey:UserID"`
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	Provider     string           `json:"provider"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type TrustedContact struct {
	gorm.Model
	UserID       uint
	ContactName  string
	ContactPhone string
	ContactEmail string
}
