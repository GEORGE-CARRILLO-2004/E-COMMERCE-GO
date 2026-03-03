package customer

import (
	"time"

	"golan/domain/customer/errors"
	"golan/domain/customer/vo"
)

type Customer struct {
	id        vo.CustomerID
	email     vo.Email
	name      string
	phone     string
	address   vo.Address
	password  vo.Password
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}

func NewCustomer(emailStr, name, phone, plainPassword string, address vo.Address) (*Customer, error) {
	if name == "" {
		return nil, errors.ErrInvalidName
	}
	if phone == "" {
		return nil, errors.ErrInvalidPhone
	}
	email, err := vo.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}
	password, err := vo.NewPassword(plainPassword)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Customer{
		id:        vo.NewCustomerID(),
		email:     email,
		name:      name,
		phone:     phone,
		address:   address,
		password:  password,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Reconstitute(
	id vo.CustomerID,
	email vo.Email,
	name string,
	phone string,
	address vo.Address,
	password vo.Password,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Customer {
	return &Customer{
		id:        id,
		email:     email,
		name:      name,
		phone:     phone,
		address:   address,
		password:  password,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *Customer) ID() vo.CustomerID     { return c.id }
func (c *Customer) Email() vo.Email       { return c.email }
func (c *Customer) Name() string          { return c.name }
func (c *Customer) Phone() string         { return c.phone }
func (c *Customer) Address() vo.Address   { return c.address }
func (c *Customer) Password() vo.Password { return c.password }
func (c *Customer) IsActive() bool        { return c.isActive }
func (c *Customer) CreatedAt() time.Time  { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time  { return c.updatedAt }

func (c *Customer) Authenticate(plainPassword string) error {
	if !c.isActive {
		return errors.ErrInactiveCustomer
	}
	return c.password.Compare(plainPassword)
}

func (c *Customer) CanPerformAction() error {
	if !c.isActive {
		return errors.ErrInactiveCustomer
	}
	return nil
}

func (c *Customer) IsOwner(id vo.CustomerID) bool {
	return c.id.String() == id.String()
}

func (c *Customer) Activate() {
	c.isActive = true
	c.updatedAt = time.Now()
}

func (c *Customer) Deactivate() {
	c.isActive = false
	c.updatedAt = time.Now()
}

func (c *Customer) UpdateName(name string) error {
	if err := c.CanPerformAction(); err != nil {
		return errors.ErrCannotUpdateInactive
	}
	if name == "" {
		return errors.ErrInvalidName
	}
	c.name = name
	c.updatedAt = time.Now()
	return nil
}

func (c *Customer) UpdatePhone(phone string) error {
	if err := c.CanPerformAction(); err != nil {
		return errors.ErrCannotUpdateInactive
	}
	if phone == "" {
		return errors.ErrInvalidPhone
	}
	c.phone = phone
	c.updatedAt = time.Now()
	return nil
}

func (c *Customer) UpdateAddress(address vo.Address) error {
	if err := c.CanPerformAction(); err != nil {
		return errors.ErrCannotUpdateInactive
	}
	c.address = address
	c.updatedAt = time.Now()
	return nil
}

func (c *Customer) UpdateEmail(emailStr string) error {
	if err := c.CanPerformAction(); err != nil {
		return errors.ErrCannotUpdateInactive
	}
	email, err := vo.NewEmail(emailStr)
	if err != nil {
		return err
	}
	c.email = email
	c.updatedAt = time.Now()
	return nil
}

func (c *Customer) ChangePassword(oldPlain, newPlain string) error {
	if err := c.CanPerformAction(); err != nil {
		return errors.ErrCannotUpdateInactive
	}
	if err := c.password.Compare(oldPlain); err != nil {
		return errors.ErrInvalidAuth
	}
	newPassword, err := vo.NewPassword(newPlain)
	if err != nil {
		return err
	}
	c.password = newPassword
	c.updatedAt = time.Now()
	return nil
}
