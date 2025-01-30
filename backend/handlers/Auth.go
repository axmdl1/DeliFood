package handlers

import (
	"DeliFood/backend/models"
	"DeliFood/backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "auth.html", nil)
		return
	}

	var form struct {
		Email         string `form:"email" binding:"required"`
		Username      string `form:"username" binding:"required"`
		Password      string `form:"password" binding:"required"`
		CheckPassword string `form:"checkPassword" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	if form.Password != form.CheckPassword {
		log.Println("Password does not match")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check if email or username exists
	exists, err := userRepo.CheckEmailOrUsernameExists(form.Email, form.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate verification code
	verificationCode := utils.GenerateVerificationCode()

	// Create user object
	user := models.User{
		UserName:         form.Username,
		Email:            form.Email,
		Password:         string(hashedPassword),
		VerificationCode: verificationCode,
		IsVerified:       false,
		Role:             "user", // Default role is user
	}

	// Save user
	if err := userRepo.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	// Send verification email
	if err := utils.SendVerificationEmail(form.Email, verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	// Redirect to verification page
	c.Redirect(http.StatusSeeOther, "/auth/verify-email?email="+form.Email)
}

// VerifyEmailHandler processes email verification
func VerifyEmailHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		email := c.Query("email")

		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is missing"})
			return
		}

		c.HTML(http.StatusOK, "verify.html", gin.H{"Email": email})
		return
	}

	var form struct {
		Email            string `form:"email" binding:"required"`
		VerificationCode string `form:"code" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email or verification code"})
		return
	}

	// Fetch user from DB
	user, err := userRepo.GetUserByEmail(form.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify code
	if user.VerificationCode != form.VerificationCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Update verification status
	if err := userRepo.VerifyEmail(form.Email, form.VerificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating verification status"})
		return
	}

	// Redirect based on user role
	if user.Role == "admin" {
		c.Redirect(http.StatusSeeOther, "/admin/panel")
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

// LoginHandler processes user login
func LoginHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "auth.html", nil)
		return
	}

	// Get email and password from the form
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Authenticate the user using the provided email and password
	user, err := userRepo.Authenticate(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed: " + err.Error()})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the token as an HTTP-only cookie (so it’s not accessible via JavaScript)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(12 * time.Hour), // Token expires in 24 hours
		HttpOnly: true,                           // Makes the cookie inaccessible via JavaScript
		Secure:   true,                           // For HTTPS, use true in production (false if not using HTTPS)
		Path:     "/",
	})

	// Redirect to the admin panel or home page
	if user.Role == "admin" {
		c.Redirect(http.StatusSeeOther, "/admin/panel")
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}
