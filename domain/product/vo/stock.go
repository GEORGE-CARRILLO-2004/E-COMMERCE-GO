package vo

import "golan/domain/product/errors"

type Stock struct {
	quantity int
}

func NewStock(quantity int) (Stock, error) {
	if quantity < 0 {
		return Stock{}, errors.ErrInvalidProductData
	}
	return Stock{quantity: quantity}, nil
}

func ReconstituteStock(quantity int) Stock {
	return Stock{quantity: quantity}
}

func (s Stock) Quantity() int {
	return s.quantity
}

func (s Stock) Decrease(qty Quantity) (Stock, error) {
	if s.quantity < qty.Value() {
		return s, errors.ErrInsufficientStock
	}
	return Stock{quantity: s.quantity - qty.Value()}, nil
}

func (s Stock) Increase(qty Quantity) Stock {
	return Stock{quantity: s.quantity + qty.Value()}
}
