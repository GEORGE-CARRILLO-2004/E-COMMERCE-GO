package review

import (
	"time"

	customerVO "golan/domain/customer/vo"
	productVO "golan/domain/product/vo"
	"golan/domain/review/errors"
	"golan/domain/review/vo"
)

type Review struct {
	id        vo.ReviewID
	productID productVO.ProductID
	authorID  customerVO.CustomerID
	rating    vo.Rating
	comment   string
	createdAt time.Time
	updatedAt time.Time
}

func NewReview(productID productVO.ProductID, authorID customerVO.CustomerID, rating vo.Rating, comment string) (*Review, error) {
	if comment == "" {
		return nil, errors.ErrInvalidReviewData
	}
	now := time.Now()
	return &Review{
		id:        vo.NewReviewID(),
		productID: productID,
		authorID:  authorID,
		rating:    rating,
		comment:   comment,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Reconstitute(
	id vo.ReviewID,
	productID productVO.ProductID,
	authorID customerVO.CustomerID,
	rating vo.Rating,
	comment string,
	createdAt time.Time,
	updatedAt time.Time,
) *Review {
	return &Review{
		id:        id,
		productID: productID,
		authorID:  authorID,
		rating:    rating,
		comment:   comment,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (r *Review) ID() vo.ReviewID                 { return r.id }
func (r *Review) ProductID() productVO.ProductID  { return r.productID }
func (r *Review) AuthorID() customerVO.CustomerID { return r.authorID }
func (r *Review) Rating() vo.Rating               { return r.rating }
func (r *Review) Comment() string                 { return r.comment }
func (r *Review) CreatedAt() time.Time            { return r.createdAt }
func (r *Review) UpdatedAt() time.Time            { return r.updatedAt }

func (r *Review) IsOwnedBy(customerID customerVO.CustomerID) bool {
	return r.authorID.String() == customerID.String()
}

func (r *Review) guardOwnership(callerID customerVO.CustomerID) error {
	if !r.IsOwnedBy(callerID) {
		return errors.ErrNotReviewOwner
	}
	return nil
}

func (r *Review) UpdateComment(newComment string, callerID customerVO.CustomerID) error {
	if err := r.guardOwnership(callerID); err != nil {
		return err
	}
	if newComment == "" {
		return errors.ErrInvalidReviewData
	}
	r.comment = newComment
	r.updatedAt = time.Now()
	return nil
}

func (r *Review) UpdateRating(newRating vo.Rating, callerID customerVO.CustomerID) error {
	if err := r.guardOwnership(callerID); err != nil {
		return err
	}
	r.rating = newRating
	r.updatedAt = time.Now()
	return nil
}
