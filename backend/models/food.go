package models

type Food struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

var Foods = []Food{
	{Name: "Burger", Category: "Burger", Image: "assets/images/menu/blog-1.jpg", Description: "Indulge in our classic Burger, featuring a juicy beef patty, crisp lettuce, ripe tomatoes, and creamy cheddar cheese, all nestled in a freshly toasted sesame bun. A timeless delight that’s perfect for any craving.", Price: "290.99"},
	{Name: "Gourmet Cheeseburger", Category: "Burger", Image: "assets/images/menu/blog-2.jpg", Description: "Savor the rich flavors of our Gourmet Cheeseburger, crafted with premium beef, caramelized onions, melted cheddar, and a touch of tangy sauce. A gourmet twist on a beloved classic", Price: "310.99"},
	{Name: "Chicken Biryani", Category: "Biryani", Image: "assets/images/menu/blog-3.jpg", Description: "Experience the aromatic delight of our Chicken Biryani—a fragrant medley of perfectly spiced basmati rice and tender chicken, slow-cooked to perfection and bursting with traditional flavors", Price: "423.99"},
	{Name: "Vegetable Biryani", Category: "Biryani", Image: "assets/images/menu/blog-4.jpg", Description: "Delight in our Vegetable Biryani, a healthy and flavorful dish packed with garden-fresh vegetables, fragrant basmati rice, and a blend of aromatic spices for a satisfying and nutritious meal.", Price: "210.99"},
	{Name: "Grilled Chicken", Category: "Chicken", Image: "assets/images/menu/blog-5.jpg", Description: "Enjoy the smoky, savory goodness of our Grilled Chicken—perfectly marinated and char-grilled to lock in the juices. A simple yet irresistible dish that’s high in flavor and protein.", Price: "390.99"},
	{Name: "Roasted Chicken", Category: "Chicken", Image: "assets/images/menu/blog-6.jpg", Description: "Treat yourself to our succulent Roasted Chicken, slow-cooked to perfection with a blend of herbs and spices, delivering a crispy golden skin and tender, flavorful meat.", Price: "412.99"},
}

func GetMenuItemCount() int {
	return len(Foods)
}

func GetPaginatedMenuItems(start, end int) []Food {
	if start >= len(Foods) {
		return []Food{}
	}
	if end > len(Foods) {
		end = len(Foods)
	}
	return Foods[start:end]
}
