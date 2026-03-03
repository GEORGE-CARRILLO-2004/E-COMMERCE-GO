package in

import (
	"context"
)

type CustomerUseCase interface {
	Register(ctx context.Context, email, name, phone, password, street, city, country string) error
	Login(ctx context.Context, email, password string) (string, error)
	UpdateProfile(ctx context.Context, customerID, name, phone, street, city, country string) error
	ChangePassword(ctx context.Context, customerID, oldPassword, newPassword string) error
}

type ProductUseCase interface {
	CreateProduct(ctx context.Context, sellerID, name, description string, priceAmt float64, stock int, category string) (string, error)
	UpdateProduct(ctx context.Context, callerID, productID, name, description string, priceAmt float64) error
	DeactivateProduct(ctx context.Context, callerID, productID string) error
	ListActiveProducts(ctx context.Context) ([]ProductDTO, error)
}

type CartUseCase interface {
	AddItem(ctx context.Context, customerID, productID string, quantity int) error
	RemoveItem(ctx context.Context, customerID, productID string) error
	ClearCart(ctx context.Context, customerID string) error
	GetCart(ctx context.Context, customerID string) (CartDTO, error)
}

type OrderUseCase interface {
	CreateOrderFromCart(ctx context.Context, customerID, street, city, country string) (string, error)
	Checkout(ctx context.Context, customerID, orderID, paymentMethod string) error
	CancelOrder(ctx context.Context, customerID, orderID string) error
	GetMyOrders(ctx context.Context, customerID string) ([]OrderDTO, error)
}

type ReviewUseCase interface {
	CreateReview(ctx context.Context, customerID, productID string, rating int, comment string) error
	UpdateReview(ctx context.Context, customerID, reviewID string, rating int, comment string) error
	DeleteReview(ctx context.Context, customerID, reviewID string) error
	GetProductReviews(ctx context.Context, productID string) ([]ReviewDTO, error)
}

type ProductDTO struct {
	ID          string  `json:"id"`
	SellerID    string  `json:"seller_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

type CartDTO struct {
	ID         string        `json:"id"`
	CustomerID string        `json:"customer_id"`
	Items      []CartItemDTO `json:"items"`
	Total      float64       `json:"total"`
	ExpiresAt  string        `json:"expires_at"`
}

type CartItemDTO struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderDTO struct {
	ID              string         `json:"id"`
	Status          string         `json:"status"`
	Total           float64        `json:"total"`
	ShippingStreet  string         `json:"shipping_street"`
	ShippingCity    string         `json:"shipping_city"`
	ShippingCountry string         `json:"shipping_country"`
	CreatedAt       string         `json:"created_at"`
	Items           []OrderItemDTO `json:"items"`
}

type OrderItemDTO struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type ReviewDTO struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	AuthorID  string `json:"author_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
