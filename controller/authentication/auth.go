package authentication

import (
	"Diploma/config"
	"Diploma/users"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"time"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Phone  string `json:"phone"`
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func setupCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// 🔥 Разрешаем только определенные домены
	allowedOrigins := map[string]bool{
		"http://localhost:3000":        true,
		"https://safetap.onrender.com": true,
	}

	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true") // 🔥 Должно быть true для cookies и headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	// ✅ Обрабатываем preflight (OPTIONS) запрос
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r)

	if r.Method == "OPTIONS" { // ✅ Preflight-запросы
		w.WriteHeader(http.StatusOK)
		return
	}

	var input struct {
		Name            string `json:"name"`
		Phone           string `json:"phone"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	// Декодируем JSON-запрос
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверяем, что пароль и подтверждающий пароль совпадают
	if input.Password != input.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли номер телефона в базе данных
	var existingUser users.User
	if err := config.DB.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
		http.Error(w, "Phone number already registered", http.StatusConflict)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	user := users.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		Provider: "local",
	}

	// Сохраняем пользователя в базе данных
	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User registered with ID: %d\n", user.ID)

	// Генерируем JWT-токен
	tokenString, err := generateToken(user.ID, user.Phone)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Сохраняем токен в базе данных
	if err := config.DB.Model(&user).Update("access_token", tokenString).Error; err != nil {
		http.Error(w, "Error saving access token", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ с токеном
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Login(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r) // ✅ Добавлено

	if r.Method == "OPTIONS" { // ✅ Разрешаем preflight-запросы
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var input struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user users.User
	if err := config.DB.Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateToken(user.ID, user.Phone)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// ✅ Устанавливаем HttpOnly Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func generateToken(userID uint, phone string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	tokenStr = strings.TrimPrefix(strings.TrimSpace(tokenStr), "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var user users.User
	if err := config.DB.Where("email = ? AND provider = ?", claims.Phone, "local").First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}

	if err := config.DB.Model(&users.User{}).Where("id = ?", claims.UserID).Update("access_token", "").Error; err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}


func getUserIdHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	var user users.User
	if err := config.DB.First(&user, uint(userID)).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}