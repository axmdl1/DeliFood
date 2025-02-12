package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/pkg/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Insert the food item into the database
	err := adminRepo.AddFood(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add food: %s", err)})
		return
	}

	// Redirect back to the admin panel or show a success message
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

func DeleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID value", http.StatusBadRequest)
		return
	}

	err = userRepo.DeleteFood(parsedID)
	if err != nil {
		http.Error(w, "Failed to delete food: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
