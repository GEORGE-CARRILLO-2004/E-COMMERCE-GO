package usecases

import (
	"context"

	"golan/application/ports/in"
	"golan/application/ports/out"
	"golan/domain/customer"
	customerVO "golan/domain/customer/vo"
)

type customerService struct {
	customerRepo  out.CustomerRepository
	tokenProvider out.TokenProvider
}

func NewCustomerService(repo out.CustomerRepository, tokenProvider out.TokenProvider) in.CustomerUseCase {
	return &customerService{
		customerRepo:  repo,
		tokenProvider: tokenProvider,
	}
}

func (s *customerService) Register(ctx context.Context, email, name, phone, password, street, city, country string) error {
	address, err := customerVO.NewAddress(street, city, country)
	if err != nil {
		return err
	}

	c, err := customer.NewCustomer(email, name, phone, password, address)
	if err != nil {
		return err
	}

	return s.customerRepo.Save(ctx, c)
}

func (s *customerService) Login(ctx context.Context, emailStr, password string) (string, error) {
	email, err := customerVO.NewEmail(emailStr)
	if err != nil {
		return "", err
	}

	c, err := s.customerRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := c.Authenticate(password); err != nil {
		return "", err
	}

	return s.tokenProvider.GenerateToken(ctx, c)
}

func (s *customerService) UpdateProfile(ctx context.Context, customerIDStr, name, phone, street, city, country string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)
	c, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	if err := c.CanPerformAction(); err != nil {
		return err
	}

	if name != "" {
		if err := c.UpdateName(name); err != nil {
			return err
		}
	}
	if phone != "" {
		if err := c.UpdatePhone(phone); err != nil {
			return err
		}
	}
	if street != "" || city != "" || country != "" {
		address, err := customerVO.NewAddress(street, city, country)
		if err != nil {
			return err
		}
		if err := c.UpdateAddress(address); err != nil {
			return err
		}
	}

	return s.customerRepo.Save(ctx, c)
}

func (s *customerService) ChangePassword(ctx context.Context, customerIDStr, oldPassword, newPassword string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)
	c, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	if err := c.ChangePassword(oldPassword, newPassword); err != nil {
		return err
	}

	return s.customerRepo.Save(ctx, c)
}
