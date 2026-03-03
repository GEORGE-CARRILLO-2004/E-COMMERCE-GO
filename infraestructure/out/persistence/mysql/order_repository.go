package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golan/application/ports/out"
	customerVO "golan/domain/customer/vo"
	"golan/domain/order"
	orderErrors "golan/domain/order/errors"
	orderVO "golan/domain/order/vo"
	productVO "golan/domain/product/vo"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) out.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Save(ctx context.Context, o *order.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO orders (id, customer_id, total, status, shipping_street, shipping_city, shipping_country, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		status=VALUES(status), shipping_street=VALUES(shipping_street),
		shipping_city=VALUES(shipping_city), shipping_country=VALUES(shipping_country),
		updated_at=VALUES(updated_at)
	`
	_, err = tx.ExecContext(ctx, query,
		o.ID().String(),
		o.CustomerID().String(),
		o.Total().Amount(),
		o.Status().String(),
		o.ShippingAddress().Street(),
		o.ShippingAddress().City(),
		o.ShippingAddress().Country(),
		o.CreatedAt(),
		o.UpdatedAt(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	itemsQuery := `INSERT IGNORE INTO order_items (order_id, product_id, name, price, quantity, subtotal) VALUES (?, ?, ?, ?, ?, ?)`
	for _, item := range o.Items() {
		_, err = tx.ExecContext(ctx, itemsQuery,
			o.ID().String(),
			item.ProductID().String(),
			item.Name(),
			item.Price().Amount(),
			item.Quantity().Value(),
			item.Subtotal().Amount(),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) FindByID(ctx context.Context, id orderVO.OrderID) (*order.Order, error) {
	query := `SELECT id, customer_id, total, status, shipping_street, shipping_city, shipping_country, created_at, updated_at FROM orders WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var oIDStr, cIDStr, statusStr, shStreet, shCity, shCountry string
	var total float64
	var createdAt, updatedAt time.Time

	err := row.Scan(&oIDStr, &cIDStr, &total, &statusStr, &shStreet, &shCity, &shCountry, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, orderErrors.ErrOrderNotFound
		}
		return nil, err
	}

	items, err := r.findOrderItems(ctx, oIDStr)
	if err != nil {
		return nil, err
	}

	return order.Reconstitute(
		id,
		customerVO.ReconstituteCustomerID(cIDStr),
		items,
		productVO.ReconstituteMoney(total),
		orderVO.ReconstituteOrderStatus(statusStr),
		orderVO.ReconstituteShippingAddress(shStreet, shCity, shCountry),
		createdAt,
		updatedAt,
	), nil
}

func (r *orderRepository) FindByCustomerID(ctx context.Context, customerID customerVO.CustomerID) ([]*order.Order, error) {
	query := `SELECT id, total, status, shipping_street, shipping_city, shipping_country, created_at, updated_at FROM orders WHERE customer_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, customerID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var oIDStr, statusStr, shStreet, shCity, shCountry string
		var total float64
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&oIDStr, &total, &statusStr, &shStreet, &shCity, &shCountry, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		items, err := r.findOrderItems(ctx, oIDStr)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order.Reconstitute(
			orderVO.ReconstituteOrderID(oIDStr),
			customerID,
			items,
			productVO.ReconstituteMoney(total),
			orderVO.ReconstituteOrderStatus(statusStr),
			orderVO.ReconstituteShippingAddress(shStreet, shCity, shCountry),
			createdAt,
			updatedAt,
		))
	}

	return orders, nil
}

func (r *orderRepository) findOrderItems(ctx context.Context, orderID string) ([]*order.OrderItem, error) {
	itemsQuery := `SELECT product_id, name, price, quantity, subtotal FROM order_items WHERE order_id = ?`
	rows, err := r.db.QueryContext(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*order.OrderItem
	for rows.Next() {
		var pID, name string
		var price, subtotal float64
		var qty int

		if err := rows.Scan(&pID, &name, &price, &qty, &subtotal); err != nil {
			return nil, err
		}

		items = append(items, order.ReconstituteOrderItem(
			productVO.ReconstituteProductID(pID),
			name,
			productVO.ReconstituteMoney(price),
			productVO.ReconstituteQuantity(qty),
			productVO.ReconstituteMoney(subtotal),
		))
	}
	return items, nil
}
