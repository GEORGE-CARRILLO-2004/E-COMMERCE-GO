package vo

import "golan/domain/product/errors"

type Quantity struct {
	value int
}

func NewQuantity(val int) (Quantity, error) {
	if val <= 0 {
		return Quantity{}, errors.ErrInvalidProductData
	}
	return Quantity{value: val}, nil
}

func ReconstituteQuantity(val int) Quantity {
	return Quantity{value: val}
}

func (q Quantity) Value() int {
	return q.value
}
