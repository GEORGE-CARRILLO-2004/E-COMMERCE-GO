package vo

import (
	"regexp"

	"golan/domain/customer/errors"
)

type Email struct {
	value string
}

func NewEmail(val string) (Email, error) {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !regex.MatchString(val) {
		return Email{}, errors.ErrInvalidEmail
	}
	return Email{value: val}, nil
}

func ReconstituteEmail(val string) Email {
	return Email{value: val}
}

func (e Email) String() string {
	return e.value
}
