package payment

import (
	"time"

	orderVO "golan/domain/order/vo"
	"golan/domain/payment/errors"
	"golan/domain/payment/vo"
	productVO "golan/domain/product/vo"
)

type Payment struct {
	id            vo.PaymentID
	orderID       orderVO.OrderID
	amount        productVO.Money
	method        vo.PaymentMethod
	status        vo.PaymentStatus
	transactionID string
	processedAt   *time.Time
	failureReason string
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPayment(orderID orderVO.OrderID, amount productVO.Money, method vo.PaymentMethod) (*Payment, error) {
	if amount.Amount() <= 0 {
		return nil, errors.ErrInvalidPaymentAmount
	}

	now := time.Now()
	return &Payment{
		id:            vo.NewPaymentID(),
		orderID:       orderID,
		amount:        amount,
		method:        method,
		status:        vo.PendingStatus(),
		transactionID: "",
		processedAt:   nil,
		failureReason: "",
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func Reconstitute(
	id vo.PaymentID,
	orderID orderVO.OrderID,
	amount productVO.Money,
	method vo.PaymentMethod,
	status vo.PaymentStatus,
	transactionID string,
	processedAt *time.Time,
	failureReason string,
	createdAt time.Time,
	updatedAt time.Time,
) *Payment {
	return &Payment{
		id:            id,
		orderID:       orderID,
		amount:        amount,
		method:        method,
		status:        status,
		transactionID: transactionID,
		processedAt:   processedAt,
		failureReason: failureReason,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (p *Payment) ID() vo.PaymentID         { return p.id }
func (p *Payment) OrderID() orderVO.OrderID { return p.orderID }
func (p *Payment) Amount() productVO.Money  { return p.amount }
func (p *Payment) Method() vo.PaymentMethod { return p.method }
func (p *Payment) Status() vo.PaymentStatus { return p.status }
func (p *Payment) TransactionID() string    { return p.transactionID }
func (p *Payment) ProcessedAt() *time.Time  { return p.processedAt }
func (p *Payment) FailureReason() string    { return p.failureReason }
func (p *Payment) CreatedAt() time.Time     { return p.createdAt }
func (p *Payment) UpdatedAt() time.Time     { return p.updatedAt }

func (p *Payment) StartProcessing(transactionID string) error {
	newStatus, err := p.status.TransitionTo("PROCESSING")
	if err != nil {
		return err
	}
	p.status = newStatus
	p.transactionID = transactionID
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) Complete() error {
	newStatus, err := p.status.TransitionTo("COMPLETED")
	if err != nil {
		return err
	}
	now := time.Now()
	p.status = newStatus
	p.processedAt = &now
	p.updatedAt = now
	return nil
}

func (p *Payment) Fail(reason string) error {
	newStatus, err := p.status.TransitionTo("FAILED")
	if err != nil {
		return err
	}
	p.status = newStatus
	p.failureReason = reason
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) Refund() error {
	newStatus, err := p.status.TransitionTo("REFUNDED")
	if err != nil {
		return err
	}
	p.status = newStatus
	p.updatedAt = time.Now()
	return nil
}
