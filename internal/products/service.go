package products

import (
	"context"

	"github.com/jackc/pgx/v5"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) GetProductById(ctx context.Context, id int64) (repo.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return repo.Product{}, ErrProductNotFound
		}
		return repo.Product{}, err
	}
	return product, nil
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	products, err := s.repo.ListProducts(ctx)

	if err != nil {
		return []repo.Product{}, err
	}

	if products == nil {
		return []repo.Product{}, nil
	}

	return products, nil
}
