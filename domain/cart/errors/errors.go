package errors

import "errors"

var (
	ErrCartExpired     = errors.New("el carrito ha expirado")
	ErrCartItemInvalid = errors.New("item de carrito invalido")
	ErrExceedsMaxStock = errors.New("la cantidad excede el stock maximo disponible")
	ErrCartEmpty       = errors.New("el carrito esta vacio")
	ErrItemNotInCart   = errors.New("el item no existe en el carrito")
	ErrNotCartOwner    = errors.New("no eres el dueño de este carrito")
)
