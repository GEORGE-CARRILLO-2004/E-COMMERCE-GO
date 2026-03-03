package errors

import "errors"

var (
	ErrInvalidReviewData    = errors.New("datos de reseña invalidos")
	ErrNotReviewOwner       = errors.New("no eres el autor de esta reseña")
	ErrReviewNotFound       = errors.New("reseña no encontrada")
	ErrMustPurchaseToReview = errors.New("debes comprar el producto para poder reseñarlo")
)
