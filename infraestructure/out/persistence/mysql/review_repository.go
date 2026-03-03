package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golan/application/ports/out"
	customerVO "golan/domain/customer/vo"
	productVO "golan/domain/product/vo"
	"golan/domain/review"
	reviewErrors "golan/domain/review/errors"
	reviewVO "golan/domain/review/vo"
)

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) out.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Save(ctx context.Context, rev *review.Review) error {
	query := `
		INSERT INTO reviews (id, product_id, author_id, rating, comment, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		rating=VALUES(rating), comment=VALUES(comment), updated_at=VALUES(updated_at)
	`
	_, err := r.db.ExecContext(ctx, query,
		rev.ID().String(),
		rev.ProductID().String(),
		rev.AuthorID().String(),
		rev.Rating().Value(),
		rev.Comment(),
		rev.CreatedAt(),
		rev.UpdatedAt(),
	)
	return err
}

func (r *reviewRepository) FindByProductID(ctx context.Context, productID productVO.ProductID) ([]*review.Review, error) {
	query := `SELECT id, product_id, author_id, rating, comment, created_at, updated_at FROM reviews WHERE product_id = ?`
	rows, err := r.db.QueryContext(ctx, query, productID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*review.Review
	for rows.Next() {
		var id, pID, aID, comment string
		var rating int
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &pID, &aID, &rating, &comment, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		reviews = append(reviews, review.Reconstitute(
			reviewVO.ReconstituteReviewID(id),
			productVO.ReconstituteProductID(pID),
			customerVO.ReconstituteCustomerID(aID),
			reviewVO.ReconstituteRating(rating),
			comment,
			createdAt,
			updatedAt,
		))
	}
	return reviews, nil
}

func (r *reviewRepository) FindByID(ctx context.Context, id reviewVO.ReviewID) (*review.Review, error) {
	query := `SELECT id, product_id, author_id, rating, comment, created_at, updated_at FROM reviews WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var rID, pID, aID, comment string
	var rating int
	var createdAt, updatedAt time.Time

	err := row.Scan(&rID, &pID, &aID, &rating, &comment, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reviewErrors.ErrReviewNotFound
		}
		return nil, err
	}

	return review.Reconstitute(
		reviewVO.ReconstituteReviewID(rID),
		productVO.ReconstituteProductID(pID),
		customerVO.ReconstituteCustomerID(aID),
		reviewVO.ReconstituteRating(rating),
		comment,
		createdAt,
		updatedAt,
	), nil
}

func (r *reviewRepository) Delete(ctx context.Context, id reviewVO.ReviewID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = ?`, id.String())
	return err
}
