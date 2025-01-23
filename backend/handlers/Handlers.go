package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/pkg/repo"
	"database/sql"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"

	_ "github.com/lib/pq"
)

var dbs *sql.DB

func SetDB(db *sql.DB) {
	dbs = db
}

// Add UserRepository as a global variable
var userRepo *repo.UserRepo

// SetUserRepo sets the user repository instance
func SetUserRepo(r *repo.UserRepo) {
	userRepo = r
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	t.ExecuteTemplate(w, "index.html", nil)
}

var tmplFuncs = template.FuncMap{
	"add": func(x, y int) int { return x + y },
	"sub": func(x, y int) int { return x - y },
	"iter": func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i + 1
		}
		return result
	},
}

const itemsPerPage = 12

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering, sorting, and pagination
	category := r.URL.Query().Get("category")
	sortParam := r.URL.Query().Get("sort")
	pageQuery := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = 1
	}

	// Filter and sort menu items
	filterSort := getMenuItems(category, sortParam)

	// Calculate total items and validate pagination indices
	totalItems := len(filterSort) // Use filtered and sorted items count
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage

	if start > totalItems { // If start index exceeds total items, return no items
		start = totalItems
	}
	if end > totalItems { // If end index exceeds total items, cap it at total items
		end = totalItems
	}
	if start > end { // Ensure start is always less than or equal to end
		start = end
	}

	paginatedItems := filterSort[start:end] // Safe slicing

	// Prepare data for the template
	data := struct {
		Items       []models.Food
		CurrentPage int
		TotalPages  int
	}{
		Items:       paginatedItems,
		CurrentPage: page,
		TotalPages:  (totalItems + itemsPerPage - 1) / itemsPerPage, // Total pages calculation
	}

	// Parse and execute the template
	tmpl := template.Must(template.New("menu.html").Funcs(tmplFuncs).ParseFiles("./frontend/menu.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Filter and sort menu items
func getMenuItems(category, sortParam string) []models.Food {
	filteredItems := []models.Food{}

	// Filter by category
	for _, item := range models.Foods {
		if category == "" || item.Category == category {
			filteredItems = append(filteredItems, item)
		}
	}

	// Sort by specified parameter
	switch sortParam {
	case "price-asc":
		sort.Slice(filteredItems, func(i, j int) bool {
			return filteredItems[i].Price < filteredItems[j].Price
		})
	case "price-desc":
		sort.Slice(filteredItems, func(i, j int) bool {
			return filteredItems[i].Price > filteredItems[j].Price
		})
	case "name":
		sort.Slice(filteredItems, func(i, j int) bool {
			return filteredItems[i].Name < filteredItems[j].Name
		})
	}

	return filteredItems
}

func ContactUsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // Allow files up to 10MB
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to process form data %v", err), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		if name == "" || email == "" || subject == "" || message == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		tempDir := "temp"
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			err = os.Mkdir(tempDir, os.ModePerm)
			if err != nil {
				http.Error(w, "Failed to create directory for file storage", http.StatusInternalServerError)
				return
			}
		}

		file, header, err := r.FormFile("attachment")
		var filePath string
		if err == nil {
			defer file.Close()

			// Save the file to the server temporarily
			filePath = filepath.Join("temp", header.Filename)
			tempFile, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Failed to save file", http.StatusInternalServerError)
				return
			}
			defer tempFile.Close()

			// Write the uploaded file to the server
			_, err = tempFile.ReadFrom(file)
			if err != nil {
				http.Error(w, "Failed to process file", http.StatusInternalServerError)
				return
			}
		}

		mail := gomail.NewMessage()
		mail.SetHeader("From", "mr.akhmedali@bk.ru")
		mail.SetHeader("To", "mr.akhmedali@bk.ru")
		mail.SetHeader("Subject", fmt.Sprintf("Contact Us: %s", subject))
		mail.SetHeader("Reply-To", email)
		mail.SetBody("text/plain", fmt.Sprintf("From: %s\nEmail: %s\nMessage: %s", name, email, message))

		if filePath != "" {
			fmt.Println("Attaching file", filePath)
			mail.Attach(filePath)
		}

		dialer := gomail.NewDialer("smtp.mail.ru", 587, "mr.akhmedali@bk.ru", "LVWZUunmUvMW8giSXLe0")
		if err := dialer.DialAndSend(mail); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send email %v", err), http.StatusInternalServerError)
			return
		}

		if filePath != "" {
			os.Remove(filePath)
		}

		fmt.Fprint(w, "Email sent successfully!")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if username == "" || email == "" || password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	/*var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}*/

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Use the SignUp method to save the user in the database
	verificationCode := repo.GenerateVerificationCode()
	user := models.User{
		UserName:         username,
		Email:            email,
		Password:         string(hashedPassword),
		VerificationCode: verificationCode,
	}

	err = userRepo.SignUp(user)
	if err != nil {
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send verification email
	err = sendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		http.Error(w, "Error sending verification email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func sendVerificationEmail(toEmail, verificationCode string) error {
	fmt.Printf("Sending email to: %s with code: %s\n", toEmail, verificationCode)

	mail := gomail.NewMessage()
	mail.SetHeader("From", "mr.akhmedali@bk.ru")
	mail.SetHeader("To", toEmail)
	mail.SetHeader("Subject", "Email Verification Code")
	mail.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", verificationCode))

	dialer := gomail.NewDialer("smtp.mail.ru", 587, "mr.akhmedali@bk.ru", "LVWZUunmUvMW8giSXLe0")
	return dialer.DialAndSend(mail)
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	if email == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Verify the code in the database
	var dbCode string
	var isVerified bool
	err := dbs.QueryRow(`
		SELECT "verificationcode", "isverified" FROM users WHERE "email" = $1`,
		email).Scan(&dbCode, &isVerified)
	if err != nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	if isVerified {
		http.Error(w, "Email already verified", http.StatusBadRequest)
		return
	}

	if dbCode != code {
		http.Error(w, "Invalid verification code", http.StatusBadRequest)
		return
	}

	// Mark the user as verified
	_, err = dbs.Exec(`
		UPDATE users SET "isverified" = $1 WHERE "email" = $2`,
		true, email)
	if err != nil {
		http.Error(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Email verified successfully!"))
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("Error fetching user: %v\n", err) // Debugging log
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("Error validating password: %v\n", err) // Debugging log
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err) // Debugging log
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Store token in the database
	err = userRepo.UpdateUserToken(user.ID, token)
	if err != nil {
		fmt.Printf("Error saving token: %v\n", err) // Debugging log
		http.Error(w, "Error saving token", http.StatusInternalServerError)
		return
	}

	// Return token to the user
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
