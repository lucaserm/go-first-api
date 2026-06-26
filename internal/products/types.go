package products

import (
	"context"

	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	GetProductById(ctx context.Context, id int64) (repo.Product, error)
	ListProducts(ctx context.Context) ([]repo.Product, error)
}
