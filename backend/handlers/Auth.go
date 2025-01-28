package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"net/url"
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
		checkPassword := r.FormValue("checkPassword")

		if checkPassword != password {
			log.Println("Password does not match")
			http.Error(w, "Password does not match", http.StatusUnauthorized)
		} else {

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

			// Redirect to the verification page
			http.Redirect(w, r, "/auth/verify-email?email="+url.QueryEscape(email), http.StatusSeeOther)
			return
		}
	}

	// If method is GET, render the register page
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./frontend/auth.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// VerifyEmailHandler renders verify.html or processes verification
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	// If the method is GET, render the verification page
	if r.Method == http.MethodGet {
		// Fetch the email from the query parameter
		email := r.URL.Query().Get("email")

		// If email is missing, show an error
		if email == "" {
			http.Error(w, "Email is missing", http.StatusBadRequest)
			return
		}

		// Pass the email to the template to prefill the form (optional)
		tmpl := template.Must(template.ParseFiles("./frontend/verify.html"))
		err := tmpl.Execute(w, map[string]string{
			"Email": email,
		})
		if err != nil {
			http.Error(w, "Failed to render verification page: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If the method is POST, process the email verification
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		email := r.FormValue("email")
		verificationCode := r.FormValue("code")

		// Validate the form fields
		if email == "" || verificationCode == "" {
			http.Error(w, "Missing email or verification code", http.StatusBadRequest)
			return
		}

		// Fetch the user from the database
		user, err := userRepo.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Check if the provided verification code matches the one stored
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

		// Check user role and redirect accordingly
		if user.Role == "admin" {
			http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		}
		return
	}

	// Handle other HTTP methods
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
