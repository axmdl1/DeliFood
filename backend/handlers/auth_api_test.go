package handlers

import (
	"DeliFood/backend/pkg/repo"
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

// Инициализация userRepo перед тестами
func setupTestDB() {
	// Подключение к тестовой базе данных
	db, err := sql.Open("postgres", "user=postgres password=yourpassword dbname=test_db sslmode=disable")
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Создание тестового userRepo
	userRepo = &repo.UserRepo{DB: db}
}

func TestRegisterHandler(t *testing.T) {
	setupTestDB() // 👈 Вызов перед тестом

	formData := `email=test@example.com&username=testuser&password=test123&checkPassword=test123`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 See Other, got %d", rr.Code)
	}
}
