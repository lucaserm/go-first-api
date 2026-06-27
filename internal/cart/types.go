package cart

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// AddItemPayload is the request body for adding a variant to the cart. The
// quantity is added to any existing quantity for the same variant.
type AddItemPayload struct {
	VariantID int64 `json:"variantId" validate:"required,gt=0"`
	Quantity  int32 `json:"quantity" validate:"required,gt=0"`
}

// UpdateItemPayload is the request body for setting the absolute quantity of an
// existing cart line.
type UpdateItemPayload struct {
	Quantity int32 `json:"quantity" validate:"required,gt=0"`
}

// CartLineItem is a single resolved cart line with variant and product details.
type CartLineItem struct {
	VariantID        int64  `json:"variantId"`
	SKU              string `json:"sku"`
	ProductName      string `json:"productName"`
	Quantity         int32  `json:"quantity"`
	UnitPriceInCents int64  `json:"unitPriceInCents"`
	LineTotalInCents int64  `json:"lineTotalInCents"`
}

// CartResponse is the full cart view returned to the authenticated user.
type CartResponse struct {
	Items           []CartLineItem `json:"items"`
	SubtotalInCents int64          `json:"subtotalInCents"`
	ItemCount       int            `json:"itemCount"`
}

type Service interface {
	GetCart(ctx context.Context, userID pgtype.UUID) (CartResponse, error)
	AddItem(ctx context.Context, userID pgtype.UUID, payload AddItemPayload) (CartResponse, error)
	UpdateItem(ctx context.Context, userID pgtype.UUID, variantID int64, payload UpdateItemPayload) (CartResponse, error)
	RemoveItem(ctx context.Context, userID pgtype.UUID, variantID int64) error
	ClearCart(ctx context.Context, userID pgtype.UUID) error
}
