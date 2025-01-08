package models

type Food struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Price       string `json:"price"`
	OldPrice    string `json:"old_price"`
}

var Foods = []Food{
	{Name: "Burger", Image: "images/menu/blog-1.jpg", Description: "A juicy beef burger...", Price: "290.99", OldPrice: "350.99"},
	{Name: "Gourmet Cheeseburger", Image: "images/menu/blog-2.jpg", Description: "Savory grilled cheeseburger...", Price: "290.99", OldPrice: "350.99"},
	{Name: "Chicken Biryani", Image: "images/menu/blog-3.jpg", Description: "Fragrant medley of basmati rice...", Price: "290.99", OldPrice: "350.99"},
	{Name: "Vegetable Biryani", Image: "images/menu/blog-4.jpg", Description: "Healthy and flavorful dish...", Price: "290.99", OldPrice: "350.99"},
	{Name: "Grilled Chicken", Image: "images/menu/blog-5.jpg", Description: "Succulent grilled chicken...", Price: "290.99", OldPrice: "350.99"},
	{Name: "Roasted Chicken", Image: "images/menu/blog-6.jpg", Description: "Whole roasted chicken...", Price: "290.99", OldPrice: "350.99"},
}
