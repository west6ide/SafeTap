package users

import "time"

type CrimeReport struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `json:"type"`
	Article    string    `json:"article"`
	Severity   string    `json:"severity"`
	Region     string    `json:"region"`
	Street     string    `json:"street"`
	House      string    `json:"house"`
	PlaceType  string    `json:"place_type"`
	Target     string    `json:"target"`
	Department string    `json:"department"`
	CrimeDate  time.Time `json:"crime_date"`
	KUSINumber string    `json:"kusi_number"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	CreatedAt  time.Time `json:"createdAt"`
}



