package order

import (
	"time"

	customerVO "golan/domain/customer/vo"
	"golan/domain/order/errors"
	"golan/domain/order/vo"
	productVO "golan/domain/product/vo"
)

type Order struct {
	id              vo.OrderID
	customerID      customerVO.CustomerID
	items           []*OrderItem
	total           productVO.Money
	status          vo.OrderStatus
	shippingAddress vo.ShippingAddress
	createdAt       time.Time
	updatedAt       time.Time
}

func NewOrder(customerID customerVO.CustomerID, items []*OrderItem, shippingAddress vo.ShippingAddress) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.ErrOrderWithoutItems
	}

	total := productVO.ZeroMoney()
	for _, item := range items {
		total = total.Add(item.Subtotal())
	}

	now := time.Now()
	return &Order{
		id:              vo.NewOrderID(),
		customerID:      customerID,
		items:           items,
		total:           total,
		status:          vo.PendingStatus(),
		shippingAddress: shippingAddress,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func Reconstitute(
	id vo.OrderID,
	customerID customerVO.CustomerID,
	items []*OrderItem,
	total productVO.Money,
	status vo.OrderStatus,
	shippingAddress vo.ShippingAddress,
	createdAt time.Time,
	updatedAt time.Time,
) *Order {
	return &Order{
		id:              id,
		customerID:      customerID,
		items:           items,
		total:           total,
		status:          status,
		shippingAddress: shippingAddress,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

func (o *Order) ID() vo.OrderID                      { return o.id }
func (o *Order) CustomerID() customerVO.CustomerID   { return o.customerID }
func (o *Order) ShippingAddress() vo.ShippingAddress { return o.shippingAddress }

func (o *Order) Items() []*OrderItem {
	clone := make([]*OrderItem, len(o.items))
	copy(clone, o.items)
	return clone
}

func (o *Order) Total() productVO.Money { return o.total }
func (o *Order) Status() vo.OrderStatus { return o.status }
func (o *Order) CreatedAt() time.Time   { return o.createdAt }
func (o *Order) UpdatedAt() time.Time   { return o.updatedAt }

func (o *Order) BelongsTo(customerID customerVO.CustomerID) bool {
	return o.customerID.String() == customerID.String()
}

func (o *Order) ItemCount() int {
	return len(o.items)
}

func (o *Order) HasProduct(productID productVO.ProductID) bool {
	for _, item := range o.items {
		if item.ProductID().String() == productID.String() {
			return true
		}
	}
	return false
}

func (o *Order) recalculateTotal() {
	total := productVO.ZeroMoney()
	for _, item := range o.items {
		total = total.Add(item.Subtotal())
	}
	o.total = total
}

func (o *Order) AddItem(item *OrderItem) error {
	if !o.status.CanModifyItems() {
		return errors.ErrCannotModifyItems
	}
	o.items = append(o.items, item)
	o.recalculateTotal()
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) Pay() error {
	if !o.status.CanBePaid() {
		return errors.ErrOrderNotPending
	}
	newStatus, err := o.status.TransitionTo("PAID")
	if err != nil {
		return err
	}
	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) Confirm() error {
	newStatus, err := o.status.TransitionTo("CONFIRMED")
	if err != nil {
		return err
	}
	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) Ship() error {
	newStatus, err := o.status.TransitionTo("SHIPPED")
	if err != nil {
		return err
	}
	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) Deliver() error {
	newStatus, err := o.status.TransitionTo("DELIVERED")
	if err != nil {
		return err
	}
	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.status.String() == "CANCELLED" {
		return errors.ErrOrderAlreadyCancelled
	}
	if o.status.String() == "DELIVERED" {
		return errors.ErrCannotCancelDelivered
	}
	newStatus, err := o.status.TransitionTo("CANCELLED")
	if err != nil {
		return err
	}
	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}
