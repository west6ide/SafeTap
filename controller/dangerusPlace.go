package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func HandleCreateCrime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type          string    `json:"type"`
		Article       string    `json:"article"`
		Severity      string    `json:"severity"`
		Region        string    `json:"region"`
		AddressStreet string    `json:"street"`
		AddressNumber string    `json:"house"`
		PlaceType     string    `json:"place_type"`
		Target        string    `json:"target"`
		Department    string    `json:"department"`
		CrimeDate     time.Time `json:"crime_date"`
		KUSINumber    string    `json:"kusi_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.AddressStreet == "" || req.AddressNumber == "" || req.Region == "" {
		http.Error(w, "Address fields are incomplete", http.StatusBadRequest)
		return
	}

	address := fmt.Sprintf("%s %s, %s", req.AddressStreet, req.AddressNumber, req.Region)
	lat, lng, err := geocodeAddress(address)
	if err != nil {
		http.Error(w, "Geocoding failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	crime := users.CrimeReport{
		Type:       req.Type,
		Article:    req.Article,
		Severity:   req.Severity,
		Region:     req.Region,
		Street:     req.AddressStreet,
		House:      req.AddressNumber,
		PlaceType:  req.PlaceType,
		Target:     req.Target,
		Department: req.Department,
		CrimeDate:  req.CrimeDate,
		KUSINumber: req.KUSINumber,
		Latitude:   lat,
		Longitude:  lng,
		CreatedAt:  time.Now(),
	}

	if err := config.DB.Create(&crime).Error; err != nil {
		http.Error(w, "Failed to save crime", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Crime saved"}`))
}

func HandleGetCrimes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var crimes []users.CrimeReport
	if err := config.DB.Find(&crimes).Error; err != nil {
		http.Error(w, "Failed to fetch crimes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(crimes)
}

func geocodeAddress(address string) (float64, float64, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	endpoint := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		url.QueryEscape(address), apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}
	if len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no geocoding results")
	}

	loc := result.Results[0].Geometry.Location
	return loc.Lat, loc.Lng, nil
}
