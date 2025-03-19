package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Store = sessions.NewCookieStore([]byte("something-very-secret"))
)

func InitDB() error { // Функция теперь возвращает ошибку
	dsn := os.Getenv("DATABASE_URL") // Получение URL базы данных из переменной окружения
	var err error
	DB, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		fmt.Println("Не удалось подключиться к базе данных:", err)
		return err // Возвращаем ошибку
	}
	fmt.Println("Соединение с базой данных успешно")
	return nil // Возвращаем nil, если всё прошло успешно
}

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("❌ Ошибка: JWT_SECRET не задан в переменных окружения!")
	}
	return secret
}
