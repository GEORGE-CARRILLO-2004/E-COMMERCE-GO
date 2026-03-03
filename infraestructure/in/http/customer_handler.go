package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"golan/application/ports/in"
	"golan/application/ports/out"
)

type CustomerHandler struct {
	customerService in.CustomerUseCase
	tokenProvider   out.TokenProvider
}

func NewCustomerHandler(s in.CustomerUseCase, tp out.TokenProvider) *CustomerHandler {
	return &CustomerHandler{customerService: s, tokenProvider: tp}
}

func (h *CustomerHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Street   string `json:"street"`
		City     string `json:"city"`
		Country  string `json:"country"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.customerService.Register(r.Context(), req.Email, req.Name, req.Phone, req.Password, req.Street, req.City, req.Country); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CustomerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.customerService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.customerService.UpdateProfile(r.Context(), customerID, req.Name, req.Phone, req.Street, req.City, req.Country); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CustomerHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	customerID := extractCustomerID(r, h.tokenProvider)
	if customerID == "" {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.customerService.ChangePassword(r.Context(), customerID, req.OldPassword, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func extractCustomerID(r *http.Request, tp out.TokenProvider) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	customerID, err := tp.ValidateToken(context.Background(), parts[1])
	if err != nil {
		return ""
	}
	return customerID
}
