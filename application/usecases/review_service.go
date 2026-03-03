package usecases

import (
	"context"

	"golan/application/ports/in"
	"golan/application/ports/out"
	customerVO "golan/domain/customer/vo"
	productVO "golan/domain/product/vo"
	"golan/domain/review"
	reviewErrors "golan/domain/review/errors"
	reviewVO "golan/domain/review/vo"
)

type reviewService struct {
	reviewRepo   out.ReviewRepository
	customerRepo out.CustomerRepository
	productRepo  out.ProductRepository
	orderRepo    out.OrderRepository
}

func NewReviewService(reviewRepo out.ReviewRepository, customerRepo out.CustomerRepository, productRepo out.ProductRepository, orderRepo out.OrderRepository) in.ReviewUseCase {
	return &reviewService{
		reviewRepo:   reviewRepo,
		customerRepo: customerRepo,
		productRepo:  productRepo,
		orderRepo:    orderRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, customerIDStr, productIDStr string, ratingVal int, comment string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)
	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	productID := productVO.ReconstituteProductID(productIDStr)
	_, err = s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	orders, err := s.orderRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}
	purchased := false
	for _, o := range orders {
		if o.HasProduct(productID) {
			purchased = true
			break
		}
	}
	if !purchased {
		return reviewErrors.ErrMustPurchaseToReview
	}

	rating, err := reviewVO.NewRating(ratingVal)
	if err != nil {
		return err
	}

	r, err := review.NewReview(productID, customerID, rating, comment)
	if err != nil {
		return err
	}

	return s.reviewRepo.Save(ctx, r)
}

func (s *reviewService) UpdateReview(ctx context.Context, customerIDStr, reviewIDStr string, ratingVal int, comment string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)
	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	reviewID := reviewVO.ReconstituteReviewID(reviewIDStr)
	r, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if !r.IsOwnedBy(customerID) {
		return reviewErrors.ErrNotReviewOwner
	}

	if comment != "" {
		if err := r.UpdateComment(comment, customerID); err != nil {
			return err
		}
	}

	if ratingVal > 0 {
		rating, err := reviewVO.NewRating(ratingVal)
		if err != nil {
			return err
		}
		if err := r.UpdateRating(rating, customerID); err != nil {
			return err
		}
	}

	return s.reviewRepo.Save(ctx, r)
}

func (s *reviewService) DeleteReview(ctx context.Context, customerIDStr, reviewIDStr string) error {
	customerID := customerVO.ReconstituteCustomerID(customerIDStr)
	cust, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if err := cust.CanPerformAction(); err != nil {
		return err
	}

	reviewID := reviewVO.ReconstituteReviewID(reviewIDStr)
	r, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if !r.IsOwnedBy(customerID) {
		return reviewErrors.ErrNotReviewOwner
	}

	return s.reviewRepo.Delete(ctx, reviewID)
}

func (s *reviewService) GetProductReviews(ctx context.Context, productIDStr string) ([]in.ReviewDTO, error) {
	productID := productVO.ReconstituteProductID(productIDStr)
	reviews, err := s.reviewRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	var dtos []in.ReviewDTO
	for _, r := range reviews {
		dtos = append(dtos, in.ReviewDTO{
			ID:        r.ID().String(),
			ProductID: r.ProductID().String(),
			AuthorID:  r.AuthorID().String(),
			Rating:    r.Rating().Value(),
			Comment:   r.Comment(),
		})
	}

	if dtos == nil {
		dtos = []in.ReviewDTO{}
	}

	return dtos, nil
}
