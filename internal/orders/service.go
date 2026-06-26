package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/ecom/internal/products"
)

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

func NewService(repo *repo.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, customerID pgtype.UUID, payload createOrderParams) (repo.Order, error) {
	if !customerID.Valid {
		return repo.Order{}, ErrCustomerIdIsRequired
	}

	if len(payload.Items) == 0 {
		return repo.Order{}, ErrItemsIsRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}

	qtx := s.repo.WithTx(tx)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	order, err := qtx.CreateOrder(ctx, customerID)
	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range payload.Items {
		variant, err := qtx.GetVariantByID(ctx, item.VariantID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return repo.Order{}, products.ErrProductNotFound
			}
			return repo.Order{}, err
		}

		if variant.Stock < item.Quantity {
			return repo.Order{}, products.ErrProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      order.ID,
			VariantID:    variant.ID,
			Quantity:     item.Quantity,
			PriceInCents: variant.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}

		err = qtx.DecreaseVariantStock(ctx, repo.DecreaseVariantStockParams{
			Stock: item.Quantity,
			ID:    variant.ID,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
