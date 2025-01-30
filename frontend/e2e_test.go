package frontend

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestUserRegistration(t *testing.T) {
	// Запуск WebDriver
	service, err := selenium.NewChromeDriverService("/usr/local/bin/chromedriver", 4444)
	if err != nil {
		t.Fatalf("Error starting WebDriver: %v", err)
	}
	defer service.Stop()

	// Подключение к WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer wd.Quit()

	// Открытие страницы регистрации
	err = wd.Get("http://localhost:5050/register.html")
	if err != nil {
		t.Fatalf("Failed to load register page: %v", err)
	}

	// Заполняем форму
	emailField, _ := wd.FindElement(selenium.ByID, "email")
	usernameField, _ := wd.FindElement(selenium.ByID, "username")
	passwordField, _ := wd.FindElement(selenium.ByID, "password")
	checkPasswordField, _ := wd.FindElement(selenium.ByID, "checkPassword")
	registerButton, _ := wd.FindElement(selenium.ByID, "register-btn")

	emailField.SendKeys("test@example.com")
	usernameField.SendKeys("testuser")
	passwordField.SendKeys("test123")
	checkPasswordField.SendKeys("test123")

	// Нажимаем кнопку "Зарегистрироваться"
	registerButton.Click()

	// Ждем 3 секунды
	time.Sleep(3 * time.Second)

	// Проверяем, был ли редирект на страницу подтверждения
	currentURL, _ := wd.CurrentURL()
	expectedURL := "http://localhost:5050/auth/verify-email"
	if currentURL != expectedURL {
		t.Errorf("Expected redirect to %s, but got %s", expectedURL, currentURL)
	}
}
