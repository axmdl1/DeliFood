package handlers

import (
	"DeliFood/backend/models"
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
	// Get user ID from context (authentication middleware should ensure it's set)
	userID, _ := c.Get("userID")

	// Get the food ID from the form
	foodID, err := strconv.Atoi(c.PostForm("food_id"))
	fmt.Println("FoodId: ", foodID)
	if err != nil || foodID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	// Get the quantity from the form
	quantity, err := strconv.Atoi(c.DefaultPostForm("quantity", "1"))
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// Fetch the food details from the database using food ID
	food, err := userRepo.GetFoodByID(foodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch food details: %s", err)})
		return
	}

	// Add the item to the cart
	err = cartRepo.AddItemToCart(userID.(int), foodID, quantity, food.Name, food.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add item to cart: %s", err)})
		return
	}

	// Redirect the user to the cart page after successfully adding the item
	c.Redirect(http.StatusFound, "/cart")
}

// GetCartItemsHandler retrieves and renders the cart page with the user's cart items
func GetCartItemsHandler(c *gin.Context) {
	// Get user ID from context (authentication middleware should ensure it's set)
	userID, _ := c.Get("userID")

	// Fetch cart items from the database
	cartItems, err := cartRepo.GetCartItems(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
		return
	}

	// Render the cart page with cart items
	c.HTML(http.StatusOK, "cart.html", gin.H{
		"cart_items":  cartItems,
		"total_price": calculateTotalPrice(cartItems),
	})
}

// calculateTotalPrice calculates the total price of all items in the cart
func calculateTotalPrice(cartItems []models.CartItem) float64 {
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += item.Price * float64(item.Quantity)
	}
	return totalPrice
}

// UpdateCartItemHandler updates the quantity of an item in the cart
func UpdateCartItemHandler(c *gin.Context) {
	// Get user ID from context
	userID, _ := c.Get("userID")

	// Get the cart item ID and new quantity
	itemID, _ := strconv.Atoi(c.Param("item_id"))
	var request struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Update the quantity in the database
	err := cartRepo.UpdateItemQuantity(userID.(int), itemID, request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

// RemoveCartItemHandler removes an item from the cart
func RemoveCartItemHandler(c *gin.Context) {
	// Get user ID from context
	userID, _ := c.Get("userID")

	// Get the cart item ID
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	// Remove the item from the cart
	err := cartRepo.RemoveItemFromCart(userID.(int), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}
