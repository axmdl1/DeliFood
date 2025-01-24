package models

import "testing"

func TestCalculateDiscount(t *testing.T) {
	price := 100
	discount := 10
	expected := 90

	result := CalculateDiscount(price, discount)
	if result != expected {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}

// Пример функции для теста
func CalculateDiscount(price int, discount int) int {
	return price - (price * discount / 100)
}
