package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFoodList(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/foods", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodList)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[\"Pizza\",\"Burger\",\"Sushi\"]`
	if rec.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

// Пример обработчика
func GetFoodList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[\"Pizza\",\"Burger\",\"Sushi\"]`))
}
