package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Mock handler для эмуляции API
func SearchProduct(w http.ResponseWriter, r *http.Request) {
	var request map[string]string
	json.NewDecoder(r.Body).Decode(&request)

	searchQuery := request["query"]

	// Пример поиска
	products := []string{"Pizza", "Burger", "Sushi"}
	var results []string
	for _, product := range products {
		if searchQuery != "" && searchQuery == product {
			results = append(results, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if len(results) > 0 {
		json.NewEncoder(w).Encode(results)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode([]string{})
	}
}

// End-to-End тест для поиска продукта
func TestEndToEndSearch(t *testing.T) {
	// Создаём запрос
	requestBody := map[string]string{"query": "Pizza"}
	jsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	// Создаём тестовый HTTP-сервер
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchProduct)
	handler.ServeHTTP(rec, req)

	// Проверяем статус код
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", rec.Code)
	}

	// Проверяем ответ
	var actualResponse []string
	if err := json.Unmarshal(rec.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	expectedResponse := []string{"Pizza"}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response %v, but got %v", expectedResponse, actualResponse)
	}
}
