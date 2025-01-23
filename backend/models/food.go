package models

type Food struct {
	ID          int     `db:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
