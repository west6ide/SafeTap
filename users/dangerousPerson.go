package users

import(
	"time"
)
type DangerousPerson struct {
    ID        uint   `gorm:"primaryKey"`
    FullName  string `json:"fullName"`
    PhotoURL  string `json:"photoUrl"`
    City      string `json:"city"`
    Address   string `json:"address"`
    CreatedAt time.Time
}
