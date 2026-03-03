package vo

import "github.com/google/uuid"

type ProductID struct {
	value string
}

func NewProductID() ProductID {
	return ProductID{value: uuid.New().String()}
}

func ReconstituteProductID(val string) ProductID {
	return ProductID{value: val}
}

func (id ProductID) String() string {
	return id.value
}
