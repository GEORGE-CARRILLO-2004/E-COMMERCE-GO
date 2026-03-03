package vo

import "github.com/google/uuid"

type CustomerID struct {
	value string
}

func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New().String()}
}

func ReconstituteCustomerID(val string) CustomerID {
	return CustomerID{value: val}
}

func (id CustomerID) String() string {
	return id.value
}
