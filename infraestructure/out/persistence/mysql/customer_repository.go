package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golan/application/ports/out"
	"golan/domain/customer"
	customerErrors "golan/domain/customer/errors"
	customerVO "golan/domain/customer/vo"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) out.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Save(ctx context.Context, c *customer.Customer) error {
	query := `
		INSERT INTO customers (id, email, name, phone, street, city, country, password_hash, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		email=VALUES(email), name=VALUES(name), phone=VALUES(phone), street=VALUES(street), city=VALUES(city), country=VALUES(country), password_hash=VALUES(password_hash), is_active=VALUES(is_active), updated_at=VALUES(updated_at)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID().String(),
		c.Email().String(),
		c.Name(),
		c.Phone(),
		c.Address().Street(),
		c.Address().City(),
		c.Address().Country(),
		c.Password().Hash(),
		c.IsActive(),
		c.CreatedAt(),
		c.UpdatedAt(),
	)
	return err
}

func (r *customerRepository) FindByID(ctx context.Context, id customerVO.CustomerID) (*customer.Customer, error) {
	query := `SELECT id, email, name, phone, street, city, country, password_hash, is_active, created_at, updated_at FROM customers WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())
	return r.scanCustomer(row)
}

func (r *customerRepository) FindByEmail(ctx context.Context, email customerVO.Email) (*customer.Customer, error) {
	query := `SELECT id, email, name, phone, street, city, country, password_hash, is_active, created_at, updated_at FROM customers WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email.String())
	return r.scanCustomer(row)
}

func (r *customerRepository) scanCustomer(row *sql.Row) (*customer.Customer, error) {
	var id, email, name, phone, street, city, country, passwordHash string
	var isActive bool
	var createdAt, updatedAt time.Time

	err := row.Scan(&id, &email, &name, &phone, &street, &city, &country, &passwordHash, &isActive, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerErrors.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer.Reconstitute(
		customerVO.ReconstituteCustomerID(id),
		customerVO.ReconstituteEmail(email),
		name,
		phone,
		customerVO.ReconstituteAddress(street, city, country),
		customerVO.ReconstitutePassword(passwordHash),
		isActive,
		createdAt,
		updatedAt,
	), nil
}
