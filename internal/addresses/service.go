package addresses

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
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

func (s *svc) ListAddresses(ctx context.Context, userID pgtype.UUID) ([]repo.Address, error) {
	addresses, err := s.repo.ListAddressesByUser(ctx, userID)
	if err != nil {
		return []repo.Address{}, err
	}

	if addresses == nil {
		return []repo.Address{}, nil
	}

	return addresses, nil
}

func (s *svc) CreateAddress(ctx context.Context, userID pgtype.UUID, payload CreateAddressPayload) (repo.Address, error) {
	params := repo.CreateAddressParams{
		UserID:        userID,
		RecipientName: payload.RecipientName,
		Line1:         payload.Line1,
		Line2:         payload.Line2,
		City:          payload.City,
		Region:        payload.Region,
		PostalCode:    payload.PostalCode,
		Country:       payload.Country,
		Phone:         payload.Phone,
		IsDefault:     payload.IsDefault,
	}

	// When the new address is the default we must clear any existing default in
	// the same transaction so exactly one default exists per user.
	if !payload.IsDefault {
		return s.repo.CreateAddress(ctx, params)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Address{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.repo.WithTx(tx)

	if err := qtx.UnsetDefaultAddressesForUser(ctx, userID); err != nil {
		return repo.Address{}, err
	}

	address, err := qtx.CreateAddress(ctx, params)
	if err != nil {
		return repo.Address{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Address{}, err
	}

	return address, nil
}

func (s *svc) UpdateAddress(ctx context.Context, userID pgtype.UUID, id int64, payload UpdateAddressPayload) (repo.Address, error) {
	params := repo.UpdateAddressForUserParams{
		ID:            id,
		UserID:        userID,
		RecipientName: payload.RecipientName,
		Line1:         payload.Line1,
		Line2:         payload.Line2,
		City:          payload.City,
		Region:        payload.Region,
		PostalCode:    payload.PostalCode,
		Country:       payload.Country,
		Phone:         payload.Phone,
		IsDefault:     payload.IsDefault,
	}

	if !payload.IsDefault {
		address, err := s.repo.UpdateAddressForUser(ctx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return repo.Address{}, ErrAddressNotFound
			}
			return repo.Address{}, err
		}
		return address, nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Address{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.repo.WithTx(tx)

	if err := qtx.UnsetDefaultAddressesForUser(ctx, userID); err != nil {
		return repo.Address{}, err
	}

	address, err := qtx.UpdateAddressForUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Address{}, ErrAddressNotFound
		}
		return repo.Address{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Address{}, err
	}

	return address, nil
}

func (s *svc) DeleteAddress(ctx context.Context, userID pgtype.UUID, id int64) error {
	rows, err := s.repo.DeleteAddressForUser(ctx, repo.DeleteAddressForUserParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrAddressNotFound
	}

	return nil
}
