package order

import (
	productVO "golan/domain/product/vo"
)

type OrderItem struct {
	productID productVO.ProductID
	name      string
	price     productVO.Money
	quantity  productVO.Quantity
	subtotal  productVO.Money
}

func NewOrderItem(
	productID productVO.ProductID,
	name string,
	price productVO.Money,
	quantity productVO.Quantity,
) *OrderItem {
	subtotal := price.Multiply(quantity.Value())
	return &OrderItem{
		productID: productID,
		name:      name,
		price:     price,
		quantity:  quantity,
		subtotal:  subtotal,
	}
}

func ReconstituteOrderItem(
	productID productVO.ProductID,
	name string,
	price productVO.Money,
	quantity productVO.Quantity,
	subtotal productVO.Money,
) *OrderItem {
	return &OrderItem{
		productID: productID,
		name:      name,
		price:     price,
		quantity:  quantity,
		subtotal:  subtotal,
	}
}

func (i *OrderItem) ProductID() productVO.ProductID { return i.productID }
func (i *OrderItem) Name() string                   { return i.name }
func (i *OrderItem) Price() productVO.Money         { return i.price }
func (i *OrderItem) Quantity() productVO.Quantity   { return i.quantity }
func (i *OrderItem) Subtotal() productVO.Money      { return i.subtotal }
