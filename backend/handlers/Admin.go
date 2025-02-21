package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/pkg/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var adminRepo *repo.AdminRepo

func SetAdminRepo(r *repo.AdminRepo) {
	adminRepo = r
}

// AdminPanelHandler renders the admin page
func AdminPanelHandler(c *gin.Context) {
	// Fetch user role and ID from the context set by AuthMiddleware
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	// Check if the user is an admin
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admins only"})
		return
	}

	foods, err := adminRepo.GetAllFoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	// Render the admin panel page
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"userID": userID,
		"role":   role,
		"Foods":  foods,
	})
}

// AddFoodHandler adds a new food item to the database
func AddFoodHandler(c *gin.Context) {
	// Parse form data from the request
	var food models.Food
	if err := c.ShouldBind(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Parse the uploaded file
	file, err := c.FormFile("image_file")
	if err != nil {
		fmt.Printf("Error retrieving file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
		return
	}

	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	uploadPath := filepath.Join("frontend", "assets", "images", "menu", uniqueFilename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file"})
		return
	}

	// Now store the unique filename in the database
	food.Image = uniqueFilename
	food.CreatedAt = time.Now()
	food.UpdatedAt = time.Now()

	if err := adminRepo.AddFood(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add food: %v", err)})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel")
}

// ChangeUserRoleHandler handles the request to change a user's role
func ChangeUserRoleHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.PostForm("userID"))
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	newRole := c.PostForm("role") // 'role' comes from the form field in the frontend
	if newRole != "user" && newRole != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	err = adminRepo.UpdateUserRole(userID, newRole) // Assuming userRepo has this method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

func DeleteFoodHandler(c *gin.Context) {
	// Retrieve the food ID from the posted form
	idStr := c.PostForm("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	// Call the repository method to delete the food
	if err := adminRepo.DeleteFood(objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete food: %v", err)})
		return
	}

	// Redirect back to the admin panel after deletion
	c.Redirect(http.StatusFound, "/admin/panel")
}
