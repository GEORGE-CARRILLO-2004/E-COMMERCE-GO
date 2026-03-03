package usecases

import (
	"context"
	"time"

	"golan/application/ports/in"
	"golan/application/ports/out"
	customerVO "golan/domain/customer/vo"
	"golan/domain/order"
	orderErrors "golan/domain/order/errors"
	orderVO "golan/domain/order/vo"
	"golan/domain/payment"
	paymentVO "golan/domain/payment/vo"
)

type orderService struct {
	orderRepo    out.OrderRepository
	cartRepo     out.CartRepository
	productRepo  out.ProductRepository
	paymentRepo  out.PaymentRepository
	customerRepo out.CustomerRepository
}

func NewOrderService(
	orderRepo out.OrderRepository,
	cartRepo out.CartRepository,
	productRepo out.ProductRepository,
	paymentRepo out.PaymentRepository,
	customerRepo out.CustomerRepository,
) in.OrderUseCase {
	return &orderService{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		productRepo:  productRepo,
		paymentRepo:  paymentRepo,
		customerRepo: customerRepo,
	}
}

func (s *orderService) CreateOrderFromCart(ctx context.Context, customerIDStr, street, city, country string) (string, error) {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return "", err
	}
	if err := cust.CanPerformAction(); err != nil {
		return "", err
	}

	c, err := s.cartRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return "", err
	}

	if !c.BelongsTo(customerID) {
		return "", orderErrors.ErrNotOrderOwner
	}

	if c.IsEmpty() {
		return "", orderErrors.ErrOrderWithoutItems
	}

	shippingAddress, err := orderVO.NewShippingAddress(street, city, country)
	if err != nil {
		return "", err
	}

	var items []*order.OrderItem
	for _, item := range c.Items() {
		p, err := s.productRepo.FindByID(ctx, item.ProductID())
		if err != nil {
			return "", err
		}

		if err := p.DecreaseStock(item.Quantity()); err != nil {
			return "", err
		}

		if err := s.productRepo.Save(ctx, p); err != nil {
			return "", err
		}

		orderItem := order.NewOrderItem(item.ProductID(), item.Name(), item.Price(), item.Quantity())
		items = append(items, orderItem)
	}

	o, err := order.NewOrder(customerID, items, shippingAddress)
	if err != nil {
		return "", err
	}

	if err := s.orderRepo.Save(ctx, o); err != nil {
		return "", err
	}

	return o.ID().String(), nil
}

func (s *orderService) Checkout(ctx context.Context, customerIDStr, orderIDStr, method string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	orderID := orderVO.ReconstituteOrderID(orderIDStr)
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !o.BelongsTo(customerID) {
		return orderErrors.ErrNotOrderOwner
	}

	paymentMethod, err := paymentVO.NewPaymentMethod(method)
	if err != nil {
		return err
	}

	pay, err := payment.NewPayment(o.ID(), o.Total(), paymentMethod)
	if err != nil {
		return err
	}

	if err := pay.StartProcessing("txn-" + o.ID().String()); err != nil {
		return err
	}

	if err := pay.Complete(); err != nil {
		return err
	}

	if err := o.Pay(); err != nil {
		return err
	}

	if err := s.paymentRepo.Save(ctx, pay); err != nil {
		return err
	}

	return s.orderRepo.Save(ctx, o)
}

func (s *orderService) CancelOrder(ctx context.Context, customerIDStr, orderIDStr string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	orderID := orderVO.ReconstituteOrderID(orderIDStr)
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !o.BelongsTo(customerID) {
		return orderErrors.ErrNotOrderOwner
	}

	if err := o.Cancel(); err != nil {
		return err
	}

	return s.orderRepo.Save(ctx, o)
}

func (s *orderService) GetMyOrders(ctx context.Context, customerIDStr string) ([]in.OrderDTO, error) {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)

	orders, err := s.orderRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	var dtos []in.OrderDTO
	for _, o := range orders {
		var itemDTOs []in.OrderItemDTO
		for _, item := range o.Items() {
			itemDTOs = append(itemDTOs, in.OrderItemDTO{
				ProductID: item.ProductID().String(),
				Name:      item.Name(),
				Price:     item.Price().Amount(),
				Quantity:  item.Quantity().Value(),
				Subtotal:  item.Subtotal().Amount(),
			})
		}
		dtos = append(dtos, in.OrderDTO{
			ID:              o.ID().String(),
			Status:          o.Status().String(),
			Total:           o.Total().Amount(),
			ShippingStreet:  o.ShippingAddress().Street(),
			ShippingCity:    o.ShippingAddress().City(),
			ShippingCountry: o.ShippingAddress().Country(),
			CreatedAt:       o.CreatedAt().Format(time.RFC3339),
			Items:           itemDTOs,
		})
	}

	if dtos == nil {
		dtos = []in.OrderDTO{}
	}

	return dtos, nil
}
