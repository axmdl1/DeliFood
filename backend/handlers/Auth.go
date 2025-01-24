package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form values
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		if email == "" || username == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Check if email or username already exists
		exists, err := userRepo.CheckEmailOrUsernameExists(email, username)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email or username already exists", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Println(string(hashedPassword))
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Generate a verification code
		verificationCode := utils.GenerateVerificationCode()

		// Save the user
		user := models.User{
			UserName:         username,
			Email:            email,
			Password:         string(hashedPassword),
			VerificationCode: verificationCode,
			IsVerified:       false,
			Role:             "user", // Default role is user
		}
		err = userRepo.Register(user)
		if err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		// Send verification email
		err = utils.SendVerificationEmail(email, verificationCode)
		if err != nil {
			http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
			return
		}

		// Redirect to login or show success message
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Registration successful! Please verify your email."))
		return
	}

	// If method is GET, render the register page
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./frontend/register.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	verificationCode := r.FormValue("code")

	if email == "" || verificationCode == "" {
		http.Error(w, "Missing email or verification code", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Compare the provided verification code with the one stored in the database
	if user.VerificationCode != verificationCode {
		http.Error(w, "Invalid verification code", http.StatusBadRequest)
		return
	}

	// Mark the user as verified
	err = userRepo.UpdateVerificationStatus(user.ID, true)
	if err != nil {
		http.Error(w, "Error updating verification status", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Your email has been verified successfully!"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Allow GET method to render login page
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./frontend/auth.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render login page: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle POST request for processing login
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Println("Login attempt for email:", email) // Log for debugging

		user, err := userRepo.Authenticate(email, password)
		if err != nil {
			fmt.Println("Login failed for email:", email, "Error:", err) // Log error
			http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"token": "%s", "role": "%s"}`, token, user.Role)))
	}

	// If the request method is not GET or POST, return a 405 Method Not Allowed error
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
