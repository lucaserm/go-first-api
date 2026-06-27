package addresses

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
)

// CreateAddressPayload is the request body for creating an address for the
// authenticated user.
type CreateAddressPayload struct {
	RecipientName string `json:"recipientName" validate:"required"`
	Line1         string `json:"line1" validate:"required"`
	Line2         string `json:"line2"`
	City          string `json:"city" validate:"required"`
	Region        string `json:"region" validate:"required"`
	PostalCode    string `json:"postalCode" validate:"required"`
	Country       string `json:"country" validate:"required,len=2"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"isDefault"`
}

// UpdateAddressPayload is the request body for replacing an existing address
// belonging to the authenticated user.
type UpdateAddressPayload struct {
	RecipientName string `json:"recipientName" validate:"required"`
	Line1         string `json:"line1" validate:"required"`
	Line2         string `json:"line2"`
	City          string `json:"city" validate:"required"`
	Region        string `json:"region" validate:"required"`
	PostalCode    string `json:"postalCode" validate:"required"`
	Country       string `json:"country" validate:"required,len=2"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"isDefault"`
}

type Service interface {
	ListAddresses(ctx context.Context, userID pgtype.UUID) ([]repo.Address, error)
	CreateAddress(ctx context.Context, userID pgtype.UUID, payload CreateAddressPayload) (repo.Address, error)
	UpdateAddress(ctx context.Context, userID pgtype.UUID, id int64, payload UpdateAddressPayload) (repo.Address, error)
	DeleteAddress(ctx context.Context, userID pgtype.UUID, id int64) error
}
