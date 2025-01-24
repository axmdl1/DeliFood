package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/pkg/repo"
	"database/sql"
	"fmt"
	"gopkg.in/gomail.v2"
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
