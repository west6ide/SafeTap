package users

import "time"

type CrimeReport struct {
	ID          uint      `gorm:"primaryKey"`
	Type        string
	Article     string
	Severity    string
	Region      string
	Street      string
	House       string
	PlaceType   string
	Target      string
	Department  string
	CrimeDate   time.Time
	KUSINumber  string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
}


