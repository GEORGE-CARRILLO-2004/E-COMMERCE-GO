package vo

import "github.com/google/uuid"

type CartID struct {
	value uuid.UUID
}

func NewCartID() CartID {
	return CartID{value: uuid.New()}
}

func ReconstituteCartID(val uuid.UUID) CartID {
	return CartID{value: val}
}

func (id CartID) Value() uuid.UUID {
	return id.value
}

func (id CartID) String() string {
	return id.value.String()
}
