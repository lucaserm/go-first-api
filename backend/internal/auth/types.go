package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type RegisterPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
}

type Service interface {
	register(ctx context.Context, payload RegisterPayload) (UserResponse, error)
	login(ctx context.Context, payload LoginPayload) (UserResponse, error)
}
