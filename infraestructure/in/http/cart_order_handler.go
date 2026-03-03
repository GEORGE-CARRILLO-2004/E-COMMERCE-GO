package http

import (
	"encoding/json"
	"net/http"

	"golan/application/ports/in"
	"golan/application/ports/out"
)

type CartOrderHandler struct {
	cartService   in.CartUseCase
	orderService  in.OrderUseCase
	tokenProvider out.TokenProvider
}

func NewCartOrderHandler(cs in.CartUseCase, os in.OrderUseCase, tp out.TokenProvider) *CartOrderHandler {
	return &CartOrderHandler{cartService: cs, orderService: os, tokenProvider: tp}
}

func (h *CartOrderHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.cartService.AddItem(r.Context(), customerID, req.ProductID, req.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CartOrderHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.cartService.RemoveItem(r.Context(), customerID, req.ProductID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CartOrderHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	if err := h.cartService.ClearCart(r.Context(), customerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CartOrderHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	cart, err := h.cartService.GetCart(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartOrderHandler) CheckoutCart(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		PaymentMethod string `json:"payment_method"`
		Street        string `json:"street"`
		City          string `json:"city"`
		Country       string `json:"country"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := h.orderService.CreateOrderFromCart(r.Context(), customerID, req.Street, req.City, req.Country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.Checkout(r.Context(), customerID, orderID, req.PaymentMethod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"order_id": orderID, "status": "completed"})
}

func (h *CartOrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetMyOrders(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
