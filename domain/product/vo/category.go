package vo

import "golan/domain/product/errors"

type Category struct {
	value string
}

func NewCategory(val string) (Category, error) {
	if val == "" {
		return Category{}, errors.ErrInvalidProductData
	}
	return Category{value: val}, nil
}

func ReconstituteCategory(val string) Category {
	return Category{value: val}
}

func (c Category) String() string {
	return c.value
}
