package vo

import "golan/domain/review/errors"

type Rating struct {
	value int
}

func NewRating(val int) (Rating, error) {
	if val < 1 || val > 5 {
		return Rating{}, errors.ErrInvalidReviewData
	}
	return Rating{value: val}, nil
}

func ReconstituteRating(val int) Rating {
	return Rating{value: val}
}

func (r Rating) Value() int {
	return r.value
}
