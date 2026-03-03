package usecases

import (
	"context"

	"golan/application/ports/in"
	"golan/application/ports/out"
	customerVO "golan/domain/customer/vo"
	"golan/domain/product"
	productErrors "golan/domain/product/errors"
	productVO "golan/domain/product/vo"
)

type productService struct {
	productRepo  out.ProductRepository
	customerRepo out.CustomerRepository
}

func NewProductService(productRepo out.ProductRepository, customerRepo out.CustomerRepository) in.ProductUseCase {
	return &productService{
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, sellerIDStr, name, desc string, priceAmt float64, stockQty int, categoryStr string) (string, error) {
	sellerID := customerVO.ReconstituteCustomerID(sellerIDStr)
	seller, err := s.customerRepo.FindByID(ctx, sellerID)
	if err != nil {
		return "", err
	}

	if err := seller.CanPerformAction(); err != nil {
		return "", err
	}

	price, err := productVO.NewMoney(priceAmt)
	if err != nil {
		return "", err
	}

	stock, err := productVO.NewStock(stockQty)
	if err != nil {
		return "", err
	}

	category, err := productVO.NewCategory(categoryStr)
	if err != nil {
		return "", err
	}

	p, err := product.NewProduct(name, desc, price, stock, category, seller.ID())
	if err != nil {
		return "", err
	}

	if err := s.productRepo.Save(ctx, p); err != nil {
		return "", err
	}

	return p.ID().String(), nil
}

func (s *productService) UpdateProduct(ctx context.Context, callerIDStr, productIDStr, name, desc string, priceAmt float64) error {
	callerID := customerVO.ReconstituteCustomerID(callerIDStr)

	caller, err := s.customerRepo.FindByID(ctx, callerID)
	if err != nil {
		return err
	}
	if err := caller.CanPerformAction(); err != nil {
		return err
	}

	productID := productVO.ReconstituteProductID(productIDStr)
	p, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	if !p.IsOwnedBy(callerID) {
		return productErrors.ErrNotProductOwner
	}

	if name != "" {
		if err := p.UpdateName(name, callerID); err != nil {
			return err
		}
	}
	if desc != "" {
		if err := p.UpdateDescription(desc, callerID); err != nil {
			return err
		}
	}
	if priceAmt > 0 {
		newPrice, err := productVO.NewMoney(priceAmt)
		if err != nil {
			return err
		}
		if err := p.UpdatePrice(newPrice, callerID); err != nil {
			return err
		}
	}

	return s.productRepo.Save(ctx, p)
}

func (s *productService) DeactivateProduct(ctx context.Context, callerIDStr, productIDStr string) error {
	callerID := customerVO.ReconstituteCustomerID(callerIDStr)

	caller, err := s.customerRepo.FindByID(ctx, callerID)
	if err != nil {
		return err
	}
	if err := caller.CanPerformAction(); err != nil {
		return err
	}

	productID := productVO.ReconstituteProductID(productIDStr)
	p, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	if err := p.Deactivate(callerID); err != nil {
		return err
	}

	return s.productRepo.Save(ctx, p)
}

func (s *productService) ListActiveProducts(ctx context.Context) ([]in.ProductDTO, error) {
	products, err := s.productRepo.FindAllActive(ctx)
	if err != nil {
		return nil, err
	}

	var dtos []in.ProductDTO
	for _, p := range products {
		dtos = append(dtos, in.ProductDTO{
			ID:          p.ID().String(),
			SellerID:    p.SellerID().String(),
			Name:        p.Name(),
			Description: p.Description(),
			Price:       p.Price().Amount(),
			Stock:       p.Stock().Quantity(),
			Category:    p.Category().String(),
		})
	}
	return dtos, nil
}
