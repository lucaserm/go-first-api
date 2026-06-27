package orders

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateOrderPayload is the checkout request body. The order is built from the
// authenticated user's cart; the client only chooses which shipping address to
// snapshot onto the order.
type CreateOrderPayload struct {
	AddressID int64 `json:"addressId" validate:"required,gt=0"`
}

// OrderLineItem is a single immutable line on an order (snapshot of the variant
// at order time).
type OrderLineItem struct {
	VariantID        int64  `json:"variantId"`
	SKU              string `json:"sku"`
	ProductName      string `json:"productName"`
	Quantity         int32  `json:"quantity"`
	UnitPriceInCents int64  `json:"unitPriceInCents"`
	LineTotalInCents int64  `json:"lineTotalInCents"`
}

// OrderResponse is the full order view returned to the authenticated customer.
type OrderResponse struct {
	ID            int64           `json:"id"`
	Status        string          `json:"status"`
	Currency      string          `json:"currency"`
	SubtotalCents int64           `json:"subtotalCents"`
	ShippingCents int64           `json:"shippingCents"`
	TaxCents      int64           `json:"taxCents"`
	TotalCents    int64           `json:"totalCents"`
	Shipping      ShippingAddress `json:"shipping"`
	Items         []OrderLineItem `json:"items"`
}

// ShippingAddress is the immutable address snapshot copied onto the order.
type ShippingAddress struct {
	RecipientName string `json:"recipientName"`
	Line1         string `json:"line1"`
	Line2         string `json:"line2"`
	City          string `json:"city"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
	Phone         string `json:"phone"`
}

type Service interface {
	PlaceOrder(ctx context.Context, customerID pgtype.UUID, payload CreateOrderPayload) (OrderResponse, error)
	GetOrder(ctx context.Context, customerID pgtype.UUID, orderID int64) (OrderResponse, error)
	ListOrders(ctx context.Context, customerID pgtype.UUID) ([]OrderResponse, error)
	UpdateStatus(ctx context.Context, customerID pgtype.UUID, orderID int64, newStatus string) (OrderResponse, error)
}
