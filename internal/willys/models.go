package willys

// Customer represents a logged-in Willys user.
type Customer struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	StoreID   string `json:"storeId"`
}

// Product represents a grocery product.
type Product struct {
	Name             string `json:"name"`
	Code             string `json:"code"`
	Price            string `json:"price"`
	PriceValue       float64 `json:"priceValue"`
	ComparePrice     string `json:"comparePrice"`
	ComparePriceUnit string `json:"comparePriceUnit"`
	Manufacturer     string `json:"manufacturer"`
	DisplayVolume    string `json:"displayVolume"`
	OutOfStock       bool   `json:"outOfStock"`
}

// CartProduct is a product in the shopping cart.
type CartProduct struct {
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	Price            string  `json:"price"`
	PriceValue       float64 `json:"priceValue"`
	Quantity         int     `json:"quantity"`
	PickQuantity     int     `json:"pickQuantity"`
	TotalPrice       string  `json:"totalPrice"`
	Manufacturer     string  `json:"manufacturer"`
	DisplayVolume    string  `json:"displayVolume"`
	ComparePrice     string  `json:"comparePrice"`
	ComparePriceUnit string  `json:"comparePriceUnit"`
}

// Cart represents the shopping cart.
type Cart struct {
	Products       []CartProduct `json:"products"`
	TotalPrice     string        `json:"totalPrice"`
	TotalItems     int           `json:"totalItems"`
	TotalUnitCount int           `json:"totalUnitCount"`
	TotalDiscount  string        `json:"totalDiscount"`
}

// Pagination holds paging info for search/browse results.
type Pagination struct {
	PageSize               int `json:"pageSize"`
	CurrentPage            int `json:"currentPage"`
	NumberOfPages          int `json:"numberOfPages"`
	TotalNumberOfResults   int `json:"totalNumberOfResults"`
}

// SearchResult is the response from product search or category browse.
type SearchResult struct {
	Results    []Product  `json:"results"`
	Pagination Pagination `json:"pagination"`
}

// Category represents a node in the category tree.
type Category struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	URL      string     `json:"url"`
	Children []Category `json:"children"`
}
