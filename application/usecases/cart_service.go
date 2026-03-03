package usecases

import (
	"context"
	"time"

	"golan/application/ports/in"
	"golan/application/ports/out"
	"golan/domain/cart"
	cartErrors "golan/domain/cart/errors"
	customerVO "golan/domain/customer/vo"
	productVO "golan/domain/product/vo"
)

type cartService struct {
	cartRepo     out.CartRepository
	productRepo  out.ProductRepository
	customerRepo out.CustomerRepository
}

func NewCartService(cartRepo out.CartRepository, productRepo out.ProductRepository, customerRepo out.CustomerRepository) in.CartUseCase {
	return &cartService{
		cartRepo:     cartRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

func (s *cartService) AddItem(ctx context.Context, customerIDStr, productIDStr string, quantity int) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	c, err := s.cartRepo.FindByCustomerID(ctx, customerID)
	if err != nil || c == nil || c.IsExpired() {
		c = cart.NewCart(customerID, 30*time.Minute)
	}

	if !c.BelongsTo(customerID) {
		return cartErrors.ErrNotCartOwner
	}

	productID := productVO.ReconstituteProductID(productIDStr)
	p, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	qty, err := productVO.NewQuantity(quantity)
	if err != nil {
		return err
	}

	item, err := cart.NewCartItem(
		p.ID(),
		p.Name(),
		p.Price(),
		qty,
		p.Stock(),
		"",
		p.Description(),
	)
	if err != nil {
		return err
	}

	if err := c.AddItem(item); err != nil {
		return err
	}

	return s.cartRepo.Save(ctx, c)
}

func (s *cartService) RemoveItem(ctx context.Context, customerIDStr, productIDStr string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	c, err := s.cartRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	if !c.BelongsTo(customerID) {
		return cartErrors.ErrNotCartOwner
	}

	productID := productVO.ReconstituteProductID(productIDStr)
	if err := c.RemoveItem(productID); err != nil {
		return err
	}

	return s.cartRepo.Save(ctx, c)
}

func (s *cartService) ClearCart(ctx context.Context, customerIDStr string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	c, err := s.cartRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	if !c.BelongsTo(customerID) {
		return cartErrors.ErrNotCartOwner
	}

	if err := c.Clear(); err != nil {
		return err
	}

	return s.cartRepo.Save(ctx, c)
}

func (s *cartService) GetCart(ctx context.Context, customerIDStr string) (in.CartDTO, error) {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)
	c, err := s.cartRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return in.CartDTO{}, err
	}

	if !c.BelongsTo(customerID) {
		return in.CartDTO{}, cartErrors.ErrNotCartOwner
	}

	var items []in.CartItemDTO
	for _, item := range c.Items() {
		items = append(items, in.CartItemDTO{
			ProductID: item.ProductID().String(),
			Name:      item.Name(),
			Price:     item.Price().Amount(),
			Quantity:  item.Quantity().Value(),
			Subtotal:  item.Subtotal().Amount(),
		})
	}

	return in.CartDTO{
		ID:         c.ID().String(),
		CustomerID: c.CustomerID().String(),
		Items:      items,
		Total:      c.Total().Amount(),
		ExpiresAt:  c.ExpiresAt().Format(time.RFC3339),
	}, nil
}
