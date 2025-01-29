package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminPanelHandler serves the admin dashboard
func AdminPanelHandler(c *gin.Context) {
	// Fetch all foods from the database (without filters and sorting)
	foods, err := userRepo.GetFood("", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load foods: " + err.Error()})
		return
	}

	// Render the admin panel template with food data
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"Foods": foods,
	})
}

/*func AddOrUpdateFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate the token in the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := userRepo.GetUserByToken(token)
	if err != nil || user.Role != "admin" {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	// Parse and save food data
	var food models.Food
	err = json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if food.ID == 0 {
		err = userRepo.AddFood(food)
	} else {
		err = userRepo.UpdateFood(food)
	}

	if err != nil {
		http.Error(w, "Failed to save food: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Food saved successfully"))
}*/

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
