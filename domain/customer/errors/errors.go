package errors

import "errors"

var (
	ErrCustomerNotFound     = errors.New("cliente no encontrado")
	ErrInactiveCustomer     = errors.New("el cliente esta inactivo")
	ErrInvalidEmail         = errors.New("formato de email invalido")
	ErrInvalidAddress       = errors.New("direccion invalida")
	ErrInvalidAuth          = errors.New("credenciales invalidas")
	ErrInvalidName          = errors.New("el nombre es obligatorio")
	ErrInvalidPhone         = errors.New("el telefono es obligatorio")
	ErrCannotUpdateInactive = errors.New("no se puede actualizar un cliente inactivo")
)
