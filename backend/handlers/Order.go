package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	// Make sure to import your cart repository package, e.g.:
	// "DeliFood/backend/pkg/repo"
)

// paymentServiceURL should point to your microservice's /pay endpoint.
var paymentServiceURL = "https://lenient-pure-muskox.ngrok-free.app/pay"

func CheckoutHandler(c *gin.Context) {
	// 1. Get the user ID from the context (set by your authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 2. Retrieve cart items for the user
	cartItems, err := cartRepo.GetCartItems(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
		return
	}

	// 3. Calculate the total price from the cart items
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += item.FoodPrice * float64(item.Quantity)
	}

	// 4. Generate a unique order ID (you can use a UUID)
	orderID := uuid.New().String()

	// 5. Create a JSON payload with orderID and totalPrice
	payload := map[string]interface{}{
		"order_id": orderID,
		"amount":   totalPrice,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payload"})
		return
	}

	// 6. Send a POST request to the payment microservice with the JSON payload
	resp, err := http.Post(paymentServiceURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to payment microservice"})
		return
	}
	defer resp.Body.Close()

	// Optional: For debugging, you can dump the response (remove in production)
	// dump, _ := httputil.DumpResponse(resp, true)
	// fmt.Println(string(dump))

	// 7. Read the response from the microservice
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from payment microservice"})
		return
	}

	// 8. Depending on your workflow, either render the HTML payment form from the microservice
	// or handle redirection. Here, we assume the microservice returns an HTML payment form.
	c.Data(http.StatusOK, "text/html; charset=utf-8", responseBody)
}
