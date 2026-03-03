package http

import (
	"encoding/json"
	"net/http"

	"golan/application/ports/in"
	"golan/application/ports/out"
)

type ReviewHandler struct {
	reviewService in.ReviewUseCase
	tokenProvider out.TokenProvider
}

func NewReviewHandler(s in.ReviewUseCase, tp out.TokenProvider) *ReviewHandler {
	return &ReviewHandler{reviewService: s, tokenProvider: tp}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID string `json:"product_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.reviewService.CreateReview(r.Context(), customerID, req.ProductID, req.Rating, req.Comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "falta el product_id", http.StatusBadRequest)
		return
	}

	reviews, err := h.reviewService.GetProductReviews(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}
