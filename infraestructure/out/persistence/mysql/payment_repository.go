package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golan/application/ports/out"
	orderVO "golan/domain/order/vo"
	"golan/domain/payment"
	paymentErrors "golan/domain/payment/errors"
	paymentVO "golan/domain/payment/vo"
	productVO "golan/domain/product/vo"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) out.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(ctx context.Context, p *payment.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, amount, method, status, transaction_id, processed_at, failure_reason, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		status=VALUES(status), transaction_id=VALUES(transaction_id), processed_at=VALUES(processed_at), failure_reason=VALUES(failure_reason), updated_at=VALUES(updated_at)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID().String(),
		p.OrderID().String(),
		p.Amount().Amount(),
		p.Method().String(),
		p.Status().String(),
		p.TransactionID(),
		p.ProcessedAt(),
		p.FailureReason(),
		p.CreatedAt(),
		p.UpdatedAt(),
	)
	return err
}

func (r *paymentRepository) FindByID(ctx context.Context, id paymentVO.PaymentID) (*payment.Payment, error) {
	query := `SELECT id, order_id, amount, method, status, transaction_id, processed_at, failure_reason, created_at, updated_at FROM payments WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var pIDStr, oIDStr, methodStr, statusStr, txID, failReason string
	var amount float64
	var processedAt *time.Time
	var createdAt, updatedAt time.Time

	err := row.Scan(&pIDStr, &oIDStr, &amount, &methodStr, &statusStr, &txID, &processedAt, &failReason, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, paymentErrors.ErrPaymentNotFound
		}
		return nil, err
	}

	return payment.Reconstitute(
		id,
		orderVO.ReconstituteOrderID(oIDStr),
		productVO.ReconstituteMoney(amount),
		paymentVO.ReconstitutePaymentMethod(methodStr),
		paymentVO.ReconstitutePaymentStatus(statusStr),
		txID,
		processedAt,
		failReason,
		createdAt,
		updatedAt,
	), nil
}
