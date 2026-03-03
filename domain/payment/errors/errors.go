package errors

import "errors"

var (
	ErrInvalidPaymentAmount = errors.New("payment amount must be greater than zero")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrInvalidPaymentState  = errors.New("invalid payment state transition")
	ErrPaymentNotFound      = errors.New("pago no encontrado")
)
