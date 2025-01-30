package handlers

import (
	"DeliFood/backend/pkg/repo"
	"DeliFood/backend/utils"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SendVerificationEmailMock(email, code string) error {
	log.Printf("Mock email sent to %s with code %s", email, code)
	return nil
}

// Load environment variables and initialize the database connection
func setupTestDB() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Construct the connection string
	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=" + dbSSLMode

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	userRepo = &repo.UserRepo{DB: db}
}

func TestRegisterHandler(t *testing.T) {
	setupTestDB()

	// Mock email function
	utils.SendVerificationEmailFunc = SendVerificationEmailMock
	defer func() { utils.SendVerificationEmailFunc = utils.SendVerificationEmail }()

	uniqueEmail := fmt.Sprintf("test%d@example.com", time.Now().UnixNano())
	formData := fmt.Sprintf("email=%s&username=testuser&password=test123&checkPassword=test123", uniqueEmail)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 See Other, got %d", rr.Code)
	}
}
