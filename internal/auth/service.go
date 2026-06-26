package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/ecom/internal/utils"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) register(ctx context.Context, payload RegisterPayload) (UserResponse, error) {
	_, err := s.repo.GetUserByEmailIgnoreCase(ctx, payload.Email)
	if err == nil {
		return UserResponse{}, ErrEmailConflict
	}

	if err != pgx.ErrNoRows {
		return UserResponse{}, err
	}

	_, err = s.repo.GetUserByUsernameIgnoreCase(ctx, payload.Username)
	if err == nil {
		return UserResponse{}, ErrUsernameConflict
	}

	if err != pgx.ErrNoRows {
		return UserResponse{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return UserResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		Username:       payload.Username,
		Email:          payload.Email,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *svc) login(ctx context.Context, payload LoginPayload) (UserResponse, error) {
	user, err := s.repo.GetUserByEmailIgnoreCase(ctx, payload.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return UserResponse{}, ErrInvalidCredentials
		}
		return UserResponse{}, err
	}

	if !utils.ComparePassword(user.HashedPassword, payload.Password) {
		return UserResponse{}, ErrInvalidCredentials
	}

	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
