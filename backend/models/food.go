package models

type Food struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Price       string `json:"price"`
}
