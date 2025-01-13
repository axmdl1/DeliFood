package handlers

import (
	"DeliFood/backend/models"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"text/template"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	t.ExecuteTemplate(w, "index.html", nil)
}

var tmplFuncs = template.FuncMap{
	"add": func(x, y int) int { return x + y },
	"sub": func(x, y int) int { return x - y },
	"iter": func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i + 1
		}
		return result
	},
}

const itemsPerPage = 3

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering, sorting, and pagination
	category := r.URL.Query().Get("category")
	sortParam := r.URL.Query().Get("sort")
	pageQuery := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = 1
	}

	// Filter and sort menu items
	filterSort := getMenuItems(category, sortParam)

	// Calculate total items and validate pagination indices
	totalItems := len(filterSort) // Use filtered and sorted items count
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage

	if start > totalItems { // If start index exceeds total items, return no items
		start = totalItems
	}
	if end > totalItems { // If end index exceeds total items, cap it at total items
		end = totalItems
	}
	if start > end { // Ensure start is always less than or equal to end
		start = end
	}

	paginatedItems := filterSort[start:end] // Safe slicing

	// Prepare data for the template
	data := struct {
		Items       []models.Food
		CurrentPage int
		TotalPages  int
	}{
		Items:       paginatedItems,
		CurrentPage: page,
		TotalPages:  (totalItems + itemsPerPage - 1) / itemsPerPage, // Total pages calculation
	}

	// Parse and execute the template
	tmpl := template.Must(template.New("menu.html").Funcs(tmplFuncs).ParseFiles("./frontend/menu.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Filter and sort menu items
func getMenuItems(category, sortParam string) []models.Food {
	filteredItems := []models.Food{}

	// Filter by category
	for _, item := range models.Foods {
		if category == "" || item.Category == category {
			filteredItems = append(filteredItems, item)
		}
	}

	// Sort by specified parameter
	switch sortParam {
	case "price-asc":
		sort.Slice(filteredItems, func(i, j int) bool {
			return filteredItems[i].Price < filteredItems[j].Price
		})
	case "price-desc":
		sort.Slice(filteredItems, func(i, j int) bool {
			return filteredItems[i].Price > filteredItems[j].Price
		})
	case "name":
		sort.Slice(filteredItems, func(i, j int) bool {
			return filteredItems[i].Name < filteredItems[j].Name
		})
	}

	return filteredItems
}
