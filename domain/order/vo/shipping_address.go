package vo

import "errors"

var ErrInvalidShippingAddress = errors.New("la dirección de envío es inválida: calle, ciudad y país son obligatorios")

type ShippingAddress struct {
	street  string
	city    string
	country string
}

func NewShippingAddress(street, city, country string) (ShippingAddress, error) {
	if street == "" || city == "" || country == "" {
		return ShippingAddress{}, ErrInvalidShippingAddress
	}
	return ShippingAddress{street: street, city: city, country: country}, nil
}

func ReconstituteShippingAddress(street, city, country string) ShippingAddress {
	return ShippingAddress{street: street, city: city, country: country}
}

func (s ShippingAddress) Street() string  { return s.street }
func (s ShippingAddress) City() string    { return s.city }
func (s ShippingAddress) Country() string { return s.country }

func (s ShippingAddress) IsEmpty() bool {
	return s.street == "" && s.city == "" && s.country == ""
}
