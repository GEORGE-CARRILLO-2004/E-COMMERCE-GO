package cart

import (
	"time"

	"golan/domain/cart/errors"
	"golan/domain/cart/vo"
	customerVO "golan/domain/customer/vo"
	productVO "golan/domain/product/vo"
)

type Cart struct {
	id         vo.CartID
	customerID customerVO.CustomerID
	items      []*CartItem
	total      productVO.Money
	expiresAt  time.Time
	createdAt  time.Time
	updatedAt  time.Time
}

func NewCart(customerID customerVO.CustomerID, expirationDuration time.Duration) *Cart {
	now := time.Now()
	return &Cart{
		id:         vo.NewCartID(),
		customerID: customerID,
		items:      []*CartItem{},
		total:      productVO.ZeroMoney(),
		expiresAt:  now.Add(expirationDuration),
		createdAt:  now,
		updatedAt:  now,
	}
}

func Reconstitute(
	id vo.CartID,
	customerID customerVO.CustomerID,
	items []*CartItem,
	total productVO.Money,
	expiresAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Cart {
	return &Cart{
		id:         id,
		customerID: customerID,
		items:      items,
		total:      total,
		expiresAt:  expiresAt,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (c *Cart) ID() vo.CartID                     { return c.id }
func (c *Cart) CustomerID() customerVO.CustomerID { return c.customerID }

func (c *Cart) Items() []*CartItem {
	clone := make([]*CartItem, len(c.items))
	copy(clone, c.items)
	return clone
}

func (c *Cart) Total() productVO.Money { return c.total }
func (c *Cart) ExpiresAt() time.Time   { return c.expiresAt }
func (c *Cart) CreatedAt() time.Time   { return c.createdAt }
func (c *Cart) UpdatedAt() time.Time   { return c.updatedAt }

func (c *Cart) BelongsTo(customerID customerVO.CustomerID) bool {
	return c.customerID.String() == customerID.String()
}

func (c *Cart) IsExpired() bool {
	return time.Now().After(c.expiresAt)
}

func (c *Cart) IsEmpty() bool {
	return len(c.items) == 0
}

func (c *Cart) ItemCount() int {
	return len(c.items)
}

func (c *Cart) recalculateTotal() {
	total := productVO.ZeroMoney()
	for _, item := range c.items {
		total = total.Add(item.Subtotal())
	}
	c.total = total
}

func (c *Cart) guardExpiration() error {
	if c.IsExpired() {
		return errors.ErrCartExpired
	}
	return nil
}

func (c *Cart) AddItem(item *CartItem) error {
	if err := c.guardExpiration(); err != nil {
		return err
	}

	for i, existing := range c.items {
		if existing.ProductID().String() == item.ProductID().String() {
			newQtyVal := existing.Quantity().Value() + item.Quantity().Value()
			newQty, err := productVO.NewQuantity(newQtyVal)
			if err != nil {
				return err
			}
			updatedItem, err := NewCartItem(
				existing.ProductID(),
				existing.Name(),
				existing.Price(),
				newQty,
				existing.MaxStock(),
				existing.ImageURL(),
				existing.Description(),
			)
			if err != nil {
				return err
			}
			c.items[i] = updatedItem
			c.recalculateTotal()
			c.updatedAt = time.Now()
			return nil
		}
	}

	c.items = append(c.items, item)
	c.recalculateTotal()
	c.updatedAt = time.Now()
	return nil
}

func (c *Cart) RemoveItem(productID productVO.ProductID) error {
	for i, item := range c.items {
		if item.ProductID().String() == productID.String() {
			c.items = append(c.items[:i], c.items[i+1:]...)
			c.recalculateTotal()
			c.updatedAt = time.Now()
			return nil
		}
	}
	return errors.ErrItemNotInCart
}

func (c *Cart) UpdateItemQuantity(productID productVO.ProductID, newQty productVO.Quantity) error {
	if err := c.guardExpiration(); err != nil {
		return err
	}

	for i, existing := range c.items {
		if existing.ProductID().String() == productID.String() {
			updatedItem, err := NewCartItem(
				existing.ProductID(),
				existing.Name(),
				existing.Price(),
				newQty,
				existing.MaxStock(),
				existing.ImageURL(),
				existing.Description(),
			)
			if err != nil {
				return err
			}
			c.items[i] = updatedItem
			c.recalculateTotal()
			c.updatedAt = time.Now()
			return nil
		}
	}
	return errors.ErrItemNotInCart
}

func (c *Cart) Clear() error {
	c.items = []*CartItem{}
	c.total = productVO.ZeroMoney()
	c.updatedAt = time.Now()
	return nil
}

func (c *Cart) HasProduct(productID productVO.ProductID) bool {
	for _, item := range c.items {
		if item.ProductID().String() == productID.String() {
			return true
		}
	}
	return false
}
