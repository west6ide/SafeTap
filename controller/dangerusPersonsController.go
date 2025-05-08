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
			City: "Каскелен",
			Address: "жилой комплекс Арнау, Проспект Абылай хана, 2/5 блок 5",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Микрорайон Алтын ауыл, 1 ​Алтын ауыл м-н",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Проспект Абылай хана, 6",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Улица Алтын кум, 12 ​Каскелен, Алматинская область",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Проспект Абылай хана, 7а",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Проспект Абылай хана, 24в",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Проспект Абылай хана, 45а",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Проспект Абылай хана, 75",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Проспект Абылай хана, 30/1",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
	}

	for _, person := range people {
		config.DB.Create(&person)
	}
}