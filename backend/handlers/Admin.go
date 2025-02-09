package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

	// Render the admin panel page
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"userID": userID,
		"role":   role,
	})
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
