package product

import (
	"time"

	customerVO "golan/domain/customer/vo"
	"golan/domain/product/errors"
	"golan/domain/product/vo"
)

type Product struct {
	id          vo.ProductID
	name        string
	description string
	price       vo.Money
	stock       vo.Stock
	category    vo.Category
	sellerID    customerVO.CustomerID
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewProduct(name, description string, price vo.Money, stock vo.Stock, category vo.Category, sellerID customerVO.CustomerID) (*Product, error) {
	if name == "" {
		return nil, errors.ErrInvalidProductData
	}
	if description == "" {
		return nil, errors.ErrDescriptionRequired
	}
	if price.Amount() <= 0 {
		return nil, errors.ErrPriceMustBePositive
	}
	now := time.Now()
	return &Product{
		id:          vo.NewProductID(),
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		category:    category,
		sellerID:    sellerID,
		isActive:    true,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func Reconstitute(
	id vo.ProductID,
	name string,
	description string,
	price vo.Money,
	stock vo.Stock,
	category vo.Category,
	sellerID customerVO.CustomerID,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		category:    category,
		sellerID:    sellerID,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (p *Product) ID() vo.ProductID                { return p.id }
func (p *Product) Name() string                    { return p.name }
func (p *Product) Description() string             { return p.description }
func (p *Product) Price() vo.Money                 { return p.price }
func (p *Product) Stock() vo.Stock                 { return p.stock }
func (p *Product) Category() vo.Category           { return p.category }
func (p *Product) SellerID() customerVO.CustomerID { return p.sellerID }
func (p *Product) IsActive() bool                  { return p.isActive }
func (p *Product) CreatedAt() time.Time            { return p.createdAt }
func (p *Product) UpdatedAt() time.Time            { return p.updatedAt }

func (p *Product) IsOwnedBy(customerID customerVO.CustomerID) bool {
	return p.sellerID.String() == customerID.String()
}

func (p *Product) mustBeOwned(callerID customerVO.CustomerID) error {
	if !p.IsOwnedBy(callerID) {
		return errors.ErrNotProductOwner
	}
	return nil
}

func (p *Product) DecreaseStock(qty vo.Quantity) error {
	if !p.isActive {
		return errors.ErrProductInactive
	}
	newStock, err := p.stock.Decrease(qty)
	if err != nil {
		return err
	}
	p.stock = newStock
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) IncreaseStock(qty vo.Quantity, callerID customerVO.CustomerID) error {
	if err := p.mustBeOwned(callerID); err != nil {
		return err
	}
	p.stock = p.stock.Increase(qty)
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) UpdateName(name string, callerID customerVO.CustomerID) error {
	if err := p.mustBeOwned(callerID); err != nil {
		return err
	}
	if name == "" {
		return errors.ErrInvalidProductData
	}
	p.name = name
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) UpdateDescription(desc string, callerID customerVO.CustomerID) error {
	if err := p.mustBeOwned(callerID); err != nil {
		return err
	}
	if desc == "" {
		return errors.ErrDescriptionRequired
	}
	p.description = desc
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) UpdatePrice(price vo.Money, callerID customerVO.CustomerID) error {
	if err := p.mustBeOwned(callerID); err != nil {
		return err
	}
	if price.Amount() <= 0 {
		return errors.ErrPriceMustBePositive
	}
	p.price = price
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) Activate(callerID customerVO.CustomerID) error {
	if err := p.mustBeOwned(callerID); err != nil {
		return err
	}
	p.isActive = true
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) Deactivate(callerID customerVO.CustomerID) error {
	if err := p.mustBeOwned(callerID); err != nil {
		return err
	}
	p.isActive = false
	p.updatedAt = time.Now()
	return nil
}
