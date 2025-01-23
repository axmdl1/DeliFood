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

	// Fetch items from the database
	foods, err := userRepo.GetFood(category, sortParam) // Use the `GetFoods` function
	if err != nil {
		http.Error(w, "Failed to load menu items", http.StatusInternalServerError)
		return
	}

	// Filter by category
	filteredFoods := []models.Food{}
	for _, food := range foods {
		if category == "" || food.Category == category {
			filteredFoods = append(filteredFoods, food)
		}
	}

	// Sort by specified parameter
	switch sortParam {
	case "price-asc":
		sort.Slice(filteredFoods, func(i, j int) bool {
			return filteredFoods[i].Price < filteredFoods[j].Price
		})
	case "price-desc":
		sort.Slice(filteredFoods, func(i, j int) bool {
			return filteredFoods[i].Price > filteredFoods[j].Price
		})
	case "name":
		sort.Slice(filteredFoods, func(i, j int) bool {
			return filteredFoods[i].Name < filteredFoods[j].Name
		})
	}

	// Calculate total items and pagination
	totalItems := len(filteredFoods)
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}
	paginatedItems := filteredFoods[start:end]

	// Prepare data for the template
	data := struct {
		Items       []models.Food
		CurrentPage int
		TotalPages  int
	}{
		Items:       paginatedItems,
		CurrentPage: page,
		TotalPages:  (totalItems + itemsPerPage - 1) / itemsPerPage,
	}

	// Parse and execute the template
	tmpl := template.Must(template.New("menu.html").Funcs(tmplFuncs).ParseFiles("./frontend/menu.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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
	// Handle GET request: Render the login page
	if r.Method == http.MethodGet {
		// Render the auth.html template
		tmpl := template.Must(template.ParseFiles("./frontend/auth.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render login page: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle POST request: Process login
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Missing email or password", http.StatusBadRequest)
			return
		}

		// Fetch user from the database
		user, err := userRepo.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Validate password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Generate a token for the user
		token, err := generateToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Store the token in the database
		err = userRepo.UpdateUserToken(user.ID, token)
		if err != nil {
			http.Error(w, "Failed to store token", http.StatusInternalServerError)
			return
		}

		// Redirect based on role (optional)
		if user.Role == "admin" {
			http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		return
	}

	// If the request method is not GET or POST, return a 405 Method Not Allowed error
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func ChangeRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	newRole := r.URL.Query().Get("role")

	if email == "" || newRole == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Update user's role in the database
	err := userRepo.UpdateUserRole(email, newRole)
	if err != nil {
		http.Error(w, "Error updating user role", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User role updated successfully"))
}

/*
func AddFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	category := r.FormValue("category")
	image := r.FormValue("image")
	description := r.FormValue("description")
	price := r.FormValue("price")

	if name == "" || category == "" || image == "" || description == "" || price == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Save food item to the database
	err := userRepo.AddFood(models.Food{
		Name:        name,
		Category:    category,
		Image:       image,
		Description: description,
		Price:       price,
	})
	if err != nil {
		http.Error(w, "Failed to add food item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Food item added successfully"))
}
*/
