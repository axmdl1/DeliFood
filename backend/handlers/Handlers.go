package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/pkg/repo"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

var foodRepo *repo.FoodRepo
var userRepo *repo.UserRepo

func SetFoodRepo(r *repo.FoodRepo) {
	foodRepo = r
}

func SetUserRepo(ur *repo.UserRepo) {
	userRepo = ur
}

// MainPageHandler serves the main index page
func MainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

const itemsPerPage = 8

// MenuHandler handles the menu listing with sorting, filtering, and pagination
func MenuHandler(c *gin.Context) {
	category := c.Query("category")
	sortParam := c.Query("sort")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	// Fetch items from the database
	foods, err := foodRepo.GetFood(category, sortParam)
	if err != nil {
		log.Printf("Error in GetFood: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load menu items"})
		return
	}

	// Filter by category
	filteredFoods := make([]models.Food, 0)
	for _, food := range foods {
		if category == "" || food.Category == category {
			filteredFoods = append(filteredFoods, food)
		}
	}

	// Sort items
	switch sortParam {
	case "price-asc":
		sort.Slice(filteredFoods, func(i, j int) bool { return filteredFoods[i].Price < filteredFoods[j].Price })
	case "price-desc":
		sort.Slice(filteredFoods, func(i, j int) bool { return filteredFoods[i].Price > filteredFoods[j].Price })
	case "name":
		sort.Slice(filteredFoods, func(i, j int) bool { return filteredFoods[i].Name < filteredFoods[j].Name })
	}

	// Pagination logic
	totalItems := len(filteredFoods)
	start := (page - 1) * itemsPerPage
	if start > totalItems {
		start = totalItems
	}
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}
	paginatedItems := filteredFoods[start:end]

	// Render the menu template
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"Items":       paginatedItems,
		"CurrentPage": page,
		"TotalPages":  (totalItems + itemsPerPage - 1) / itemsPerPage,
	})
}

// ContactUsHandler processes contact form submissions
func ContactUsHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	subject := c.PostForm("subject")
	message := c.PostForm("message")

	if name == "" || email == "" || subject == "" || message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Create temp directory for file uploads
	tempDir := "temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err = os.Mkdir(tempDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}
	}

	// Handle file upload
	file, err := c.FormFile("attachment")
	var filePath string
	if err == nil {
		filePath = filepath.Join(tempDir, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	}

	// Prepare email
	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("DELIFOOD_MAIL"))
	mail.SetHeader("To", os.Getenv("DELIFOOD_MAIL"))
	mail.SetHeader("Subject", fmt.Sprintf("Contact Us: %s", subject))
	mail.SetHeader("Reply-To", email)
	mail.SetBody("text/plain", fmt.Sprintf("From: %s\nEmail: %s\nMessage: %s", name, email, message))

	// Attach file if exists
	if filePath != "" {
		log.Println("Attaching file:", filePath)
		mail.Attach(filePath)
	}

	// Send email
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("DELIFOOD_MAIL"), os.Getenv("DELIFOOD_PASSWORD"))
	if err := dialer.DialAndSend(mail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// Remove file after sending
	if filePath != "" {
		os.Remove(filePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully!"})
}
