package commands

import (
	"context"

	"github.com/google/uuid"
)

type AddToCartCommand struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}

type AddToCartHandler struct {
	cartRepo cart.Repository
}

func NewAddToCartHandler(cartRepo cart.Repository) *AddToCartHandler {
	return &AddToCartHandler{cartRepo: cartRepo}
}

func (h *AddToCartHandler) Handle(ctx context.Context, cmd AddToCartCommand) error {
	cart, err := h.cartRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	if err := cart.AddItem(cmd.ProductID, cmd.Quantity); err != nil {
		return err
	}

	return &h.cartRepo.Save(ctx, cart)
}
