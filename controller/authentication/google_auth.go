package authentication

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
)

func init() {
	// Проверка, что все переменные окружения заданы
	if GoogleOauthConfig.ClientID == "" || GoogleOauthConfig.ClientSecret == "" || GoogleOauthConfig.RedirectURL == "" {
		log.Fatal("Не установлены переменные окружения для Google OAuth")
	}

	// Настройки для сессий (опционально для безопасности)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // Время жизни сессии в секундах (8 часов)
		HttpOnly: true,     // Куки не доступны через JavaScript
		Secure:   false,    // Используйте true для HTTPS
		SameSite: http.SameSiteStrictMode,
	}
}

// HandleGoogleLogin initiates Google OAuth login
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := "google"
	url := GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback processes the OAuth callback and retrieves user info from Google
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := "google"
	if r.FormValue("state") != state {
		log.Println("Invalid OAuth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := GoogleOauthConfig.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Printf("Error while exchanging code for token: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Error while fetching user info: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Convert JSON response to structure
	var userInfo map[string]interface{}
	if err := json.Unmarshal(content, &userInfo); err != nil {
		log.Printf("Error parsing user info: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Extract user info with type assertion
	googleID, ok := userInfo["id"].(string)
	if !ok {
		log.Println("Error extracting Google ID")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		log.Println("Error extracting email")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	firstName, ok := userInfo["given_name"].(string)
	if !ok {
		log.Println("Error extracting first name")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	lastName, ok := userInfo["family_name"].(string)
	if !ok {
		log.Println("Error extracting last name")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Проверка, существует ли пользователь с таким email
	var user users.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// Если пользователь не найден, создаем нового
		if err == gorm.ErrRecordNotFound {
			log.Printf("Пользователь с email %s не найден, создаем нового", email)
			user = users.User{
				Email:       email,
				Name:        firstName + " " + lastName,
				Provider:    "google",
				AccessToken: token.AccessToken,
			}
			if err := config.DB.Create(&user).Error; err != nil {
				log.Printf("Ошибка при создании пользователя: %v", err)
				http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Ошибка при попытке найти пользователя с email %s: %v", email, err)
			http.Error(w, "Ошибка поиска пользователя", http.StatusInternalServerError)
			return
		}
	}

	// Проверка в таблице GoogleUser
	var googleUser users.GoogleUser
	if err := config.DB.Where("google_id = ?", googleID).First(&googleUser).Error; err != nil {
		// Если GoogleUser не найден, создаем нового
		if err == gorm.ErrRecordNotFound {
			log.Printf("GoogleUser с ID %s не найден, создаем нового", googleID)
			googleUser = users.GoogleUser{
				UserID:      user.ID,
				GoogleID:    googleID,
				Email:       email,
				FirstName:   firstName,
				LastName:    lastName,
				AccessToken: token.AccessToken,
			}
			if err := config.DB.Create(&googleUser).Error; err != nil {
				log.Printf("Ошибка при создании GoogleUser: %v", err)
				http.Error(w, "Ошибка создания GoogleUser", http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Ошибка при попытке найти GoogleUser с Google ID %s: %v", googleID, err)
			http.Error(w, "Ошибка поиска GoogleUser", http.StatusInternalServerError)
			return
		}
	} else {
		// Если GoogleUser найден, обновляем информацию
		googleUser.Email = email
		googleUser.FirstName = firstName
		googleUser.LastName = lastName
		googleUser.AccessToken = token.AccessToken
		if err := config.DB.Save(&googleUser).Error; err != nil {
			log.Printf("Ошибка при обновлении GoogleUser: %v", err)
			http.Error(w, "Ошибка обновления GoogleUser", http.StatusInternalServerError)
			return
		}
	}

	// Сохраняем минимальные данные пользователя в сессии (например, только ID)
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Ошибка при получении сессии: %s", err.Error())
		http.Error(w, "Ошибка получения сессии", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = user.ID // Сохраняем только ID пользователя
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Ошибка при сохранении сессии: %s", err.Error())
		http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
		return
	}

	// Перенаправляем пользователя на защищенную страницу
	http.Redirect(w, r, "/welcome", http.StatusTemporaryRedirect)
}

// ValidateGoogleToken проверяет переданный токен и возвращает пользователя Google
func ValidateGoogleToken(r *http.Request) (*users.GoogleUser, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid Authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Проверяем, есть ли пользователь с этим токеном в БД
	var googleUser users.GoogleUser
	if err := config.DB.Where("access_token = ?", accessToken).First(&googleUser).Error; err == nil {
		log.Println("Токен найден в БД, используем кешированные данные")
		return &googleUser, nil
	}

	// Если токена нет в БД, проверяем его через Google API
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	defer resp.Body.Close()

	// Парсим ответ
	var tokenInfo struct {
		Sub       string `json:"sub"`
		Email     string `json:"email"`
		Audience  string `json:"aud"`
		ExpiresIn string `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse token info: %v", err)
	}

	// Преобразуем `expires_in` в int
	expiresIn, err := strconv.Atoi(tokenInfo.ExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_in value: %v", err)
	}

	// Проверяем, что токен принадлежит нашему клиенту
	if tokenInfo.Audience != os.Getenv("GOOGLE_CLIENT_ID") {
		return nil, errors.New("invalid token audience")
	}

	// Проверяем, существует ли пользователь в БД
	if err := config.DB.Where("google_id = ?", tokenInfo.Sub).First(&googleUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found in the database")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Обновляем токен, если он изменился
	if googleUser.AccessToken != accessToken {
		googleUser.AccessToken = accessToken
		if err := config.DB.Save(&googleUser).Error; err != nil {
			return nil, fmt.Errorf("failed to update access token: %v", err)
		}
	}

	log.Printf("Токен действителен, истекает через %d секунд", expiresIn)

	return &googleUser, nil
}
