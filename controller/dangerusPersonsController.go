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
			City:     "Алматы",
			Address:  "Улица Тимирязева, 59Б ​Бостандыкский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Улица Шопена, 25 ​Бостандыкский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Асыл Тобе, жилой комплекс Улица Егизбаева, 7/21 ​Бостандыкский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Улица Шевченко, 149 ​Алмалинский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "O'NER, жилой комплекс ​ЖК ONer​улица Ислама Каримова, 203 блок 2 Алмалинский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Улица Тургут Озала, 94 ​Алмалинский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Улица Прокофьева, 41 Тастак-2 м-н, Алмалинский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Жетысу-3 микрорайон, 67 Жетысу-3 м-н, Ауэзовский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Алматы",
			Address:  "Улица Жубанова, 7 Ауэзовский район",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},







		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Микрорайон Алтын ауыл, 6 ​Алтын ауыл м-н",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Проспект Абылай хана, 26Б/2 ",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Арнау, жилой комплекс Проспект Абылай хана, 2/5 блок 5",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Проспект Абылай хана, 19",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Улица 10 лет Независимости, 56",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Улица Абен Омерали, 47 · проспект Абылай хана, 82",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Улица Карасай батыра, 2",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Алатау Ажары, жилой комплекс Аубая Байгазиева улица, 35Б блок Б",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Каскелен",
			Address:  "Улица Бауыржана Момышулы, 3",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
	}

	for _, person := range people {
		config.DB.Create(&person)
	}
}
