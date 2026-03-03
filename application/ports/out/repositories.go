package out

import (
	"context"

	"golan/domain/cart"
	cartVO "golan/domain/cart/vo"
	"golan/domain/customer"
	customerVO "golan/domain/customer/vo"
	"golan/domain/order"
	orderVO "golan/domain/order/vo"
	"golan/domain/payment"
	paymentVO "golan/domain/payment/vo"
	"golan/domain/product"
	productVO "golan/domain/product/vo"
	"golan/domain/review"
	reviewVO "golan/domain/review/vo"
)

type CustomerRepository interface {
	Save(ctx context.Context, c *customer.Customer) error
	FindByID(ctx context.Context, id customerVO.CustomerID) (*customer.Customer, error)
	FindByEmail(ctx context.Context, email customerVO.Email) (*customer.Customer, error)
}

type ProductRepository interface {
	Save(ctx context.Context, p *product.Product) error
	FindByID(ctx context.Context, id productVO.ProductID) (*product.Product, error)
	FindAllActive(ctx context.Context) ([]*product.Product, error)
}

type OrderRepository interface {
	Save(ctx context.Context, o *order.Order) error
	FindByID(ctx context.Context, id orderVO.OrderID) (*order.Order, error)
	FindByCustomerID(ctx context.Context, customerID customerVO.CustomerID) ([]*order.Order, error)
}

type PaymentRepository interface {
	Save(ctx context.Context, p *payment.Payment) error
	FindByID(ctx context.Context, id paymentVO.PaymentID) (*payment.Payment, error)
}

type ReviewRepository interface {
	Save(ctx context.Context, r *review.Review) error
	FindByID(ctx context.Context, id reviewVO.ReviewID) (*review.Review, error)
	FindByProductID(ctx context.Context, productID productVO.ProductID) ([]*review.Review, error)
	Delete(ctx context.Context, id reviewVO.ReviewID) error
}

type CartRepository interface {
	Save(ctx context.Context, c *cart.Cart) error
	FindByCustomerID(ctx context.Context, id customerVO.CustomerID) (*cart.Cart, error)
	Delete(ctx context.Context, id cartVO.CartID) error
}
