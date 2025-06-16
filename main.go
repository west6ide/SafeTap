package main

import (
	"Diploma/config"
	"Diploma/controller"
	authentication2 "Diploma/controller/authentication"
	"Diploma/users"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Инициализация базы данных
	if err := config.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	controller.StartNotificationCleaner()
	controller.SeedDangerousPeople()

	// Миграция базы данных
	if err := config.DB.AutoMigrate(
		&users.User{},
		&users.TrustedContact{},
		&users.LiveLocation{},
		&users.Notification{},
		&users.SOSSignal{},
		&users.FakeCall{},
		&users.DangerousPerson{},
		&users.CrimeReport{},
		&users.SharedRoute{},
		&controller.ContactRequest{},
		&users.GoogleUser{}); err != nil {
		log.Fatalf("Ошибка миграции БД: %v", err)
	}

	// Проверка подключения к БД
	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatalf("Ошибка получения подключения к БД: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	} else {
		log.Println("Подключение к БД успешно")
	}

	handler := http.NewServeMux()

	// Эндпоинты
	handler.HandleFunc("/", handleHome)
	handler.HandleFunc("/login/google", authentication2.HandleGoogleLogin)
	handler.HandleFunc("/callback/google", authentication2.HandleGoogleCallback)
	handler.HandleFunc("/register", authentication2.Register)
	handler.HandleFunc("/login", authentication2.Login)
	handler.HandleFunc("/profile", authentication2.GetProfile)
	handler.HandleFunc("/profile/update", controller.UpdateProfile)
	handler.HandleFunc("/logout", authentication2.Logout)

	handler.HandleFunc("/contacts/add", controller.AddEmergencyContact)       // Добавление
	handler.HandleFunc("/contacts", controller.GetEmergencyContacts)          // Получение всех контактов
	handler.HandleFunc("/contacts/delete", controller.DeleteEmergencyContact) // Удаление
	handler.HandleFunc("/contacts/request", controller.SendContactRequest) // Новая реализация
	handler.HandleFunc("/contacts/respond", controller.HandleContactRequest)

	handler.HandleFunc("/sos", controller.SaveSOS)
	handler.HandleFunc("/notifications", controller.GetNotifications)
	handler.HandleFunc("/getUserId", authentication2.GetUserIdHandler)

	handler.HandleFunc("/location/update", controller.UpdateLiveLocation)
	handler.HandleFunc("/location/emergency", controller.GetEmergencyContactsLocations)

	handler.HandleFunc("/fake-calls/create", controller.ScheduleFakeCall)
	handler.HandleFunc("/fake-calls/list", controller.GetUserFakeCalls)
	handler.HandleFunc("/fake-calls/delete", controller.DeleteFakeCall)

	handler.HandleFunc("/dangerous-people", controller.GetDangerousPeople)

	// Регистрация
	handler.HandleFunc("/crimes", controller.HandleCreateCrime)
	handler.HandleFunc("/crimes/get", controller.HandleGetCrimes)

	handler.HandleFunc("/share_route", controller.ShareRouteHandler(config.DB))
	handler.HandleFunc("/shared_routes", controller.GetSharedRoutesHandler(config.DB))

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(handler)

	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session-name")
	user := session.Values["user"]

	if user != nil {
		switch usr := user.(type) {
		case users.GoogleUser:
			html := fmt.Sprintf(`
				<html><body>
				<p>Добро пожаловать, %s!</p>
				<a href="/logout">Выйти</a><br>
				<form action="/google-logout" method="post">
					<button type="submit">Выйти из Google</button>
				</form>
				</body></html>`, usr.FirstName)
			fmt.Fprint(w, html)
		default:
			http.Error(w, "Неизвестный тип пользователя", http.StatusInternalServerError)
		}
	} else {
		html := `<html><body>
		<a href="/login/google">Войти через Google</a><br>
		</body></html>`
		fmt.Fprint(w, html)
	}
}
