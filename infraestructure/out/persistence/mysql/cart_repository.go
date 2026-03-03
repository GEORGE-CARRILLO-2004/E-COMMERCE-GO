package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"golan/application/ports/out"
	"golan/domain/cart"
	cartVO "golan/domain/cart/vo"
	customerVO "golan/domain/customer/vo"
	productVO "golan/domain/product/vo"
)

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) out.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Save(ctx context.Context, c *cart.Cart) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO carts (id, customer_id, total, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		total=VALUES(total), expires_at=VALUES(expires_at), updated_at=VALUES(updated_at)
	`
	_, err = tx.ExecContext(ctx, query,
		c.ID().String(),
		c.CustomerID().String(),
		c.Total().Amount(),
		c.ExpiresAt(),
		c.CreatedAt(),
		c.UpdatedAt(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = ?`, c.ID().String())
	if err != nil {
		tx.Rollback()
		return err
	}

	itemsQuery := `INSERT INTO cart_items (cart_id, product_id, name, price, quantity, subtotal, max_stock, image_url, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, item := range c.Items() {
		_, err = tx.ExecContext(ctx, itemsQuery,
			c.ID().String(),
			item.ProductID().String(),
			item.Name(),
			item.Price().Amount(),
			item.Quantity().Value(),
			item.Subtotal().Amount(),
			item.MaxStock().Quantity(),
			item.ImageURL(),
			item.Description(),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *cartRepository) FindByCustomerID(ctx context.Context, id customerVO.CustomerID) (*cart.Cart, error) {
	query := `SELECT id, total, expires_at, created_at, updated_at FROM carts WHERE customer_id = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var cIDStr string
	var total float64
	var expiresAt, createdAt, updatedAt time.Time

	err := row.Scan(&cIDStr, &total, &expiresAt, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	cID, err := uuid.Parse(cIDStr)
	if err != nil {
		return nil, err
	}

	itemsQuery := `SELECT product_id, name, price, quantity, subtotal, max_stock, image_url, description FROM cart_items WHERE cart_id = ?`
	rows, err := r.db.QueryContext(ctx, itemsQuery, cIDStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*cart.CartItem
	for rows.Next() {
		var pID, name, img, desc string
		var price, subtotal float64
		var qty, maxStock int

		if err := rows.Scan(&pID, &name, &price, &qty, &subtotal, &maxStock, &img, &desc); err != nil {
			return nil, err
		}

		items = append(items, cart.ReconstituteCartItem(
			productVO.ReconstituteProductID(pID),
			name,
			productVO.ReconstituteMoney(price),
			productVO.ReconstituteQuantity(qty),
			productVO.ReconstituteMoney(subtotal),
			productVO.ReconstituteStock(maxStock),
			img,
			desc,
		))
	}

	return cart.Reconstitute(
		cartVO.ReconstituteCartID(cID),
		id,
		items,
		productVO.ReconstituteMoney(total),
		expiresAt,
		createdAt,
		updatedAt,
	), nil
}

func (r *cartRepository) Delete(ctx context.Context, id cartVO.CartID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = ?`, id.String())
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM carts WHERE id = ?`, id.String())
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
