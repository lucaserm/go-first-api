package orders

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	Items []orderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, customerID pgtype.UUID, payload createOrderParams) (repo.Order, error)
}
