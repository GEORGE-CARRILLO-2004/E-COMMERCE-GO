package errors

import "errors"

var (
	ErrOrderWithoutItems     = errors.New("no se puede crear una orden sin items")
	ErrInvalidOrderState     = errors.New("transicion de estado invalida")
	ErrCannotModifyItems     = errors.New("no se pueden modificar items despues de confirmar")
	ErrOrderNotPending       = errors.New("la orden no esta en estado pendiente")
	ErrNotOrderOwner         = errors.New("no eres el dueño de esta orden")
	ErrOrderAlreadyCancelled = errors.New("la orden ya esta cancelada")
	ErrCannotCancelDelivered = errors.New("no se puede cancelar una orden entregada")
	ErrOrderNotFound         = errors.New("orden no encontrada")
)
