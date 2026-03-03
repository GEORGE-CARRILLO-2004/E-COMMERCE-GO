package errors

import "errors"

var (
	ErrInvalidProductData  = errors.New("datos de producto invalidos")
	ErrProductNotFound     = errors.New("producto no encontrado")
	ErrInsufficientStock   = errors.New("stock insuficiente")
	ErrProductInactive     = errors.New("el producto esta inactivo")
	ErrNotProductOwner     = errors.New("no eres el dueño de este producto")
	ErrPriceMustBePositive = errors.New("el precio debe ser mayor a cero")
	ErrDescriptionRequired = errors.New("la descripcion es obligatoria")
)
