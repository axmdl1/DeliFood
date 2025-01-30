package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler обрабатывает регистрацию пользователей
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Разбираем форму
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Получаем данные из формы
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		checkPassword := r.FormValue("checkPassword")

		// Проверяем совпадение паролей
		if checkPassword != password {
			log.Println("❌ Passwords do not match")
			http.Error(w, "Passwords do not match", http.StatusUnauthorized)
			return
		}

		// Проверяем, что поля не пустые
		if email == "" || username == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли уже пользователь с таким email или username
		exists, err := userRepo.CheckEmailOrUsernameExists(email, username)
		if err != nil {
			log.Println("❌ Error checking existing user:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if exists {
			log.Println("⚠️ Email or username already exists:", email, username)
			http.Error(w, "Email or username already exists", http.StatusBadRequest)
			return
		}

		// Хешируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ Error hashing password:", err)
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Генерируем код подтверждения
		verificationCode := utils.GenerateVerificationCode()

		// Создаем пользователя
		user := models.User{
			UserName:         username,
			Email:            email,
			Password:         string(hashedPassword),
			VerificationCode: verificationCode,
			IsVerified:       false,
			Role:             "user",
		}

		// Сохраняем пользователя в БД
		err = userRepo.Register(user)
		if err != nil {
			log.Println("❌ Error saving user:", err) // Логируем ошибку
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		// Отправляем email с кодом подтверждения
		err = utils.SendVerificationEmail(email, verificationCode)
		if err != nil {
			log.Println("❌ Error sending verification email:", err)
			http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
			return
		}

		// Перенаправляем пользователя на страницу подтверждения
		http.Redirect(w, r, "/auth/verify-email?email="+url.QueryEscape(email), http.StatusSeeOther)
		return
	}

	// Если метод GET, рендерим страницу регистрации
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./frontend/auth.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Println("❌ Error rendering registration page:", err)
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// VerifyEmailHandler обрабатывает подтверждение email
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email is missing", http.StatusBadRequest)
			return
		}

		tmpl := template.Must(template.ParseFiles("./frontend/verify.html"))
		err := tmpl.Execute(w, map[string]string{"Email": email})
		if err != nil {
			log.Println("❌ Error rendering verification page:", err)
			http.Error(w, "Failed to render verification page", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		verificationCode := r.FormValue("code")

		if email == "" || verificationCode == "" {
			http.Error(w, "Missing email or verification code", http.StatusBadRequest)
			return
		}

		user, err := userRepo.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.VerificationCode != verificationCode {
			http.Error(w, "Invalid verification code", http.StatusBadRequest)
			return
		}

		err = userRepo.VerifyEmail(email, verificationCode)
		if err != nil {
			http.Error(w, "Error updating verification status", http.StatusInternalServerError)
			return
		}

		if user.Role == "admin" {
			http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// LoginHandler обрабатывает вход в систему
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./frontend/auth.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Println("❌ Error rendering login page:", err)
			http.Error(w, "Failed to render login page", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Println("Login attempt for email:", email)

		user, err := userRepo.Authenticate(email, password)
		if err != nil {
			fmt.Println("Login failed for email:", email, "Error:", err)
			http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		if user.Role == "admin" {
			http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"token": "%s", "role": "%s"}`, token, user.Role)))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
