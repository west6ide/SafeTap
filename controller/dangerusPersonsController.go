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
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "ул. Абулхаир хана",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Улица Розыбакиева, 111",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "11-й микрорайон, 4а",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Микрорайон Коктем-1, 44а",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Микрорайон Казахфильм, 24 · улица Исиналиева, 24",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Улица Садвакасова, 108",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "Улица Мажорова, 21",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Алматы",
			Address: "жилой комплекс Евразия,Улица Масанчи, 23/4",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
	}

	for _, person := range people {
		config.DB.Create(&person)
	}
}