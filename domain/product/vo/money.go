package vo

import (
	"math"

	"golan/domain/product/errors"
)

type Money struct {
	amount float64
}

func NewMoney(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, errors.ErrInvalidProductData
	}
	return Money{amount: math.Round(amount*100) / 100}, nil
}

func ZeroMoney() Money {
	return Money{amount: 0}
}

func ReconstituteMoney(amount float64) Money {
	return Money{amount: amount}
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Add(other Money) Money {
	return Money{amount: math.Round((m.amount+other.amount)*100) / 100}
}

func (m Money) Multiply(quantity int) Money {
	return Money{amount: math.Round(m.amount*float64(quantity)*100) / 100}
}
