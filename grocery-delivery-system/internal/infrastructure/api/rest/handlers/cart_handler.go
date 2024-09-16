package handlers

import (
	"encoding/json"
	"grocery-delivery/internal/application/commands"
	"net/http"

	"github.com/google/uuid"
)

type CartHandler struct {
	addToCartHandler *commands.AddToCartHandler
}

func NewCartHandler(addToCartHandler *commands.AddToCartHandler) *CartHandler {
	return &CartHandler{addToCartHandler: addToCartHandler}
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    uuid.UUID `json:"user_id"`
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.AddToCartCommand{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := h.addToCartHandler.Handle(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
