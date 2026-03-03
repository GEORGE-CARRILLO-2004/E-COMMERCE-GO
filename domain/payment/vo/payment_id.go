package vo

import "github.com/google/uuid"

type PaymentID struct {
	value string
}

func NewPaymentID() PaymentID {
	return PaymentID{value: uuid.New().String()}
}

func ReconstitutePaymentID(val string) PaymentID {
	return PaymentID{value: val}
}

func (id PaymentID) String() string {
	return id.value
}
