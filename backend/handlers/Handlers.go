package handlers

import (
	"DeliFood/backend/models"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

// Struct to hold page-related data
type MenuData struct {
	Items       []models.Food
	CurrentPage int
	TotalPages  int
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		fmt.Printf(err.Error())
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
	pageQuery := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = 1
	}

	totalItems := models.GetMenuItemCount() // Replace with actual implementation
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}

	items := models.GetPaginatedMenuItems(start, end) // Replace with actual implementation

	data := MenuData{
		Items:       items,
		CurrentPage: page,
		TotalPages:  (totalItems + itemsPerPage - 1) / itemsPerPage,
	}

	tmpl := template.Must(template.New("menu.html").Funcs(tmplFuncs).ParseFiles("./frontend/menu.html"))
	tmpl.Execute(w, data)
}
