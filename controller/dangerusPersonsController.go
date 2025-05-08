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
			Address: "Улица Болашак, 5",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "11-й микрорайон, 4а",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Улица Тлендиева, 16 ",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Улица Тлендиева, 28",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Улица Курылысшы, 22",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Улица Кошкарбаева, 18",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City: "Каскелен",
			Address: "Улица Аскарова, 73",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
	}

	for _, person := range people {
		config.DB.Create(&person)
	}
}