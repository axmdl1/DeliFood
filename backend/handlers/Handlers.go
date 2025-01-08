package handlers

import (
	"DeliFood/backend/models"
	"net/http"
	"strconv"
	"text/template"
)

// Assuming you have a Foods slice that holds all food items
var Foods []models.Food

// Struct to hold page-related data
type PageData struct {
	Foods       []models.Food
	CurrentPage int
	TotalPages  int
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current page from the query parameter
	pageParam := r.URL.Query().Get("page")
	page := 1
	if pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}
	}

	// Set the number of items per page
	itemsPerPage := 10

	// Calculate the total number of pages
	totalItems := len(Foods)
	totalPages := totalItems / itemsPerPage
	if totalItems%itemsPerPage != 0 {
		totalPages++
	}

	// Calculate the start and end index for the current page
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Slice the Foods array for the current page
	pageFoods := Foods[startIndex:endIndex]

	// Prepare data to pass to the template
	data := PageData{
		Foods:       pageFoods,
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	// Render the template with the page data
	tmpl, err := template.ParseFiles("cmd/frontend/menu.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
