package vo

import "golan/domain/payment/errors"

type PaymentMethod struct {
	value string
}

func NewPaymentMethod(val string) (PaymentMethod, error) {
	if val == "" {
		return PaymentMethod{}, errors.ErrInvalidPaymentMethod
	}
	return PaymentMethod{value: val}, nil
}

func ReconstitutePaymentMethod(val string) PaymentMethod {
	return PaymentMethod{value: val}
}

func (m PaymentMethod) String() string {
	return m.value
}
