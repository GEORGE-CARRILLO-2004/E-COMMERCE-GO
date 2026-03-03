package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golan/application/ports/out"
	customerVO "golan/domain/customer/vo"
	"golan/domain/product"
	productErrors "golan/domain/product/errors"
	productVO "golan/domain/product/vo"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) out.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Save(ctx context.Context, p *product.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, category, seller_id, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		name=VALUES(name), description=VALUES(description), price=VALUES(price), stock=VALUES(stock), category=VALUES(category), is_active=VALUES(is_active), updated_at=VALUES(updated_at)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID().String(),
		p.Name(),
		p.Description(),
		p.Price().Amount(),
		p.Stock().Quantity(),
		p.Category().String(),
		p.SellerID().String(),
		p.IsActive(),
		p.CreatedAt(),
		p.UpdatedAt(),
	)
	return err
}

func (r *productRepository) FindByID(ctx context.Context, id productVO.ProductID) (*product.Product, error) {
	query := `SELECT id, name, description, price, stock, category, seller_id, is_active, created_at, updated_at FROM products WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())
	return r.scanProduct(row)
}

func (r *productRepository) FindAllActive(ctx context.Context) ([]*product.Product, error) {
	query := `SELECT id, name, description, price, stock, category, seller_id, is_active, created_at, updated_at FROM products WHERE is_active = true`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		p, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) scanProduct(row *sql.Row) (*product.Product, error) {
	var id, name, description, category, sellerID string
	var price float64
	var stock int
	var isActive bool
	var createdAt, updatedAt time.Time

	err := row.Scan(&id, &name, &description, &price, &stock, &category, &sellerID, &isActive, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, productErrors.ErrProductNotFound
		}
		return nil, err
	}

	return product.Reconstitute(
		productVO.ReconstituteProductID(id),
		name,
		description,
		productVO.ReconstituteMoney(price),
		productVO.ReconstituteStock(stock),
		productVO.ReconstituteCategory(category),
		customerVO.ReconstituteCustomerID(sellerID),
		isActive,
		createdAt,
		updatedAt,
	), nil
}

func (r *productRepository) scanProductFromRows(rows *sql.Rows) (*product.Product, error) {
	var id, name, description, category, sellerID string
	var price float64
	var stock int
	var isActive bool
	var createdAt, updatedAt time.Time

	err := rows.Scan(&id, &name, &description, &price, &stock, &category, &sellerID, &isActive, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return product.Reconstitute(
		productVO.ReconstituteProductID(id),
		name,
		description,
		productVO.ReconstituteMoney(price),
		productVO.ReconstituteStock(stock),
		productVO.ReconstituteCategory(category),
		customerVO.ReconstituteCustomerID(sellerID),
		isActive,
		createdAt,
		updatedAt,
	), nil
}
