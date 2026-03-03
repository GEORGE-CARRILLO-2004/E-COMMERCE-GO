package http

import (
	"encoding/json"
	"net/http"

	"golan/application/ports/in"
	"golan/application/ports/out"
)

type ProductHandler struct {
	productService in.ProductUseCase
	tokenProvider  out.TokenProvider
}

func NewProductHandler(s in.ProductUseCase, tp out.TokenProvider) *ProductHandler {
	return &ProductHandler{productService: s, tokenProvider: tp}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	sellerID := extractCustomerID(r, h.tokenProvider)
	if sellerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Category    string  `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.productService.CreateProduct(r.Context(), sellerID, req.Name, req.Description, req.Price, req.Stock, req.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.ListActiveProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}
