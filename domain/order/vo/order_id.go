package vo

import "github.com/google/uuid"

type OrderID struct {
	value string
}

func NewOrderID() OrderID {
	return OrderID{value: uuid.New().String()}
}

func ReconstituteOrderID(val string) OrderID {
	return OrderID{value: val}
}

func (id OrderID) String() string {
	return id.value
}
