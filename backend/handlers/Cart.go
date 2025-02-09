package handlers

import (
	"DeliFood/backend/pkg/repo"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Global variable for CartRepo
var cartRepo *repo.CartRepo

// SetCartRepo sets the cart repository for use in handlers
func SetCartRepo(cr *repo.CartRepo) {
	cartRepo = cr
}

// AddToCartHandler adds an item to the user's cart
func AddToCartHandler(c *gin.Context) {
	// Get user ID from context (from JWT)
	userID, _ := c.Get("userID")

	foodID, err := strconv.Atoi(c.DefaultQuery("food_id", "0"))
	if err != nil || foodID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	quantity, err := strconv.Atoi(c.DefaultQuery("quantity", "1"))
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// Add item to cart
	err = cartRepo.AddItemToCart(userID.(int), foodID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add item to cart: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

// GetCartHandler retrieves all items in the user's cart
func GetCartHandler(c *gin.Context) {
	// Get user ID from context
	userID, _ := c.Get("userID")

	// Get cart items for the user
	cartItems, err := cartRepo.GetCartItems(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart_items": cartItems})
}
