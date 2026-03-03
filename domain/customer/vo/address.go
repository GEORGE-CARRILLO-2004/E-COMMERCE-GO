package vo

import "golan/domain/customer/errors"

type Address struct {
	street  string
	city    string
	country string
}

func NewAddress(street, city, country string) (Address, error) {
	if street == "" || city == "" || country == "" {
		return Address{}, errors.ErrInvalidAddress
	}
	return Address{
		street:  street,
		city:    city,
		country: country,
	}, nil
}

func ReconstituteAddress(street, city, country string) Address {
	return Address{
		street:  street,
		city:    city,
		country: country,
	}
}

func (a Address) Street() string {
	return a.street
}

func (a Address) City() string {
	return a.city
}

func (a Address) Country() string {
	return a.country
}
