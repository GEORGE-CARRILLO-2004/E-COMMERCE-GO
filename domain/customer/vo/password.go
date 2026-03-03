package vo

import (
	"golan/domain/customer/errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash string
}

func NewPassword(plain string) (Password, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}
	return Password{hash: string(hashed)}, nil
}

func ReconstitutePassword(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Hash() string {
	return p.hash
}

func (p Password) Compare(plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain)); err != nil {
		return errors.ErrInvalidAuth
	}
	return nil
}
