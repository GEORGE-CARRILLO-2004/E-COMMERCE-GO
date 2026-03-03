package cart

import (
	"golan/domain/cart/errors"
	productVO "golan/domain/product/vo"
)

type CartItem struct {
	productID   productVO.ProductID
	name        string
	price       productVO.Money
	quantity    productVO.Quantity
	subtotal    productVO.Money
	maxStock    productVO.Stock
	imageURL    string
	description string
}

func NewCartItem(
	productID productVO.ProductID,
	name string,
	price productVO.Money,
	quantity productVO.Quantity,
	maxStock productVO.Stock,
	imageURL string,
	description string,
) (*CartItem, error) {
	if quantity.Value() > maxStock.Quantity() {
		return nil, errors.ErrExceedsMaxStock
	}

	subtotal := price.Multiply(quantity.Value())

	return &CartItem{
		productID:   productID,
		name:        name,
		price:       price,
		quantity:    quantity,
		subtotal:    subtotal,
		maxStock:    maxStock,
		imageURL:    imageURL,
		description: description,
	}, nil
}

func ReconstituteCartItem(
	productID productVO.ProductID,
	name string,
	price productVO.Money,
	quantity productVO.Quantity,
	subtotal productVO.Money,
	maxStock productVO.Stock,
	imageURL string,
	description string,
) *CartItem {
	return &CartItem{
		productID:   productID,
		name:        name,
		price:       price,
		quantity:    quantity,
		subtotal:    subtotal,
		maxStock:    maxStock,
		imageURL:    imageURL,
		description: description,
	}
}

func (i *CartItem) ProductID() productVO.ProductID { return i.productID }
func (i *CartItem) Name() string                   { return i.name }
func (i *CartItem) Price() productVO.Money         { return i.price }
func (i *CartItem) Quantity() productVO.Quantity   { return i.quantity }
func (i *CartItem) Subtotal() productVO.Money      { return i.subtotal }
func (i *CartItem) MaxStock() productVO.Stock      { return i.maxStock }
func (i *CartItem) ImageURL() string               { return i.imageURL }
func (i *CartItem) Description() string            { return i.description }
