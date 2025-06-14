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
			City:     "Астана",
			Address:  "Айдидар, жилой комплекс Улица Жумекен Нажимеденов, 16",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Медиполь, жилой комплекс Улица Аманжол Болекпаев, 3",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Улица Арганаты, 11",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Хайвил Астана блок Е1, жилой комплекс Проспект Ракымжан Кошкарбаев, 8",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Ак Булак-3, жилой комплекс Переулок Тасшокы, 3",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Проспект Бауыржан Момышулы, 22/4",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Olymp Palace, жилой комплекс ​ЖК Олимп палас​улица Туркестан, 8",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Sensata Park, жилой комплекс Улица Туркестан, 16/4",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
		{
			FullName: "Имя Фамилия",
			City:     "Астана",
			Address:  "Promenade Expo Block D, жилой комплекс Проспект Мангилик Ел, 51",
			PhotoURL: "https://www.mirgovorit.ru/static/users/default_profile_image.9acfe78b8e1c.png",
		},
	}

	for _, person := range people {
		config.DB.Create(&person)
	}
}
