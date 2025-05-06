package controller

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"net/http"
)

func GetDangerousPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var people []users.DangerousPerson
	if err := config.DB.Find(&people).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(people)
}

func SeedDangerousPeople() {
	people := []users.DangerousPerson{
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Микрорайон-3",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "ул. Абулхаир хана",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Улица Розыбакиева, 111",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "11-й микрорайон, 4а",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Микрорайон Коктем-1, 44а",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Микрорайон Казахфильм, 24 · улица Исиналиева, 24",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Улица Садвакасова, 108",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Улица Мажорова, 21",
			PhotoURL: "",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "жилой комплекс Евразия,Улица Масанчи, 23/4",
			PhotoURL: "",
		},
	}

	for _, person := range people {
		config.DB.Create(&person)
	}
}