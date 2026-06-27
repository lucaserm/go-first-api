package cart

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

// getOrCreateCart returns the user's cart, creating it on first use. CreateCart
// uses ON CONFLICT DO NOTHING, so a concurrent insert is handled by re-reading
// the cart when the insert returns no rows.
func (s *svc) getOrCreateCart(ctx context.Context, userID pgtype.UUID) (repo.Cart, error) {
	cart, err := s.repo.GetCartByUser(ctx, userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return repo.Cart{}, err
	}

	cart, err = s.repo.CreateCart(ctx, userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return repo.Cart{}, err
	}

	// Lost the race to a concurrent creator; the cart now exists.
	return s.repo.GetCartByUser(ctx, userID)
}

func (s *svc) GetCart(ctx context.Context, userID pgtype.UUID) (CartResponse, error) {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return CartResponse{}, err
	}

	return s.buildCartResponse(ctx, cart.ID)
}

func (s *svc) AddItem(ctx context.Context, userID pgtype.UUID, payload AddItemPayload) (CartResponse, error) {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return CartResponse{}, err
	}

	variant, err := s.repo.GetVariantByID(ctx, payload.VariantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CartResponse{}, ErrVariantNotFound
		}
		return CartResponse{}, err
	}

	existingQty, err := s.existingQuantity(ctx, cart.ID, payload.VariantID)
	if err != nil {
		return CartResponse{}, err
	}

	if int64(existingQty)+int64(payload.Quantity) > int64(variant.Stock) {
		return CartResponse{}, ErrInsufficientStock
	}

	if _, err := s.repo.UpsertCartItem(ctx, repo.UpsertCartItemParams{
		CartID:    cart.ID,
		VariantID: payload.VariantID,
		Quantity:  payload.Quantity,
	}); err != nil {
		return CartResponse{}, err
	}

	return s.buildCartResponse(ctx, cart.ID)
}

func (s *svc) UpdateItem(ctx context.Context, userID pgtype.UUID, variantID int64, payload UpdateItemPayload) (CartResponse, error) {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return CartResponse{}, err
	}

	variant, err := s.repo.GetVariantByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CartResponse{}, ErrVariantNotFound
		}
		return CartResponse{}, err
	}

	if int64(payload.Quantity) > int64(variant.Stock) {
		return CartResponse{}, ErrInsufficientStock
	}

	if _, err := s.repo.SetCartItemQuantity(ctx, repo.SetCartItemQuantityParams{
		CartID:    cart.ID,
		VariantID: variantID,
		Quantity:  payload.Quantity,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CartResponse{}, ErrCartItemNotFound
		}
		return CartResponse{}, err
	}

	return s.buildCartResponse(ctx, cart.ID)
}

func (s *svc) RemoveItem(ctx context.Context, userID pgtype.UUID, variantID int64) error {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	rows, err := s.repo.DeleteCartItem(ctx, repo.DeleteCartItemParams{
		CartID:    cart.ID,
		VariantID: variantID,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrCartItemNotFound
	}

	return nil
}

func (s *svc) ClearCart(ctx context.Context, userID pgtype.UUID) error {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	return s.repo.ClearCart(ctx, cart.ID)
}

// existingQuantity returns the quantity currently held for a variant in the
// cart, or zero when the variant is not yet present.
func (s *svc) existingQuantity(ctx context.Context, cartID int64, variantID int64) (int32, error) {
	items, err := s.repo.ListCartItemsWithVariant(ctx, cartID)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		if item.VariantID == variantID {
			return item.Quantity, nil
		}
	}

	return 0, nil
}

// buildCartResponse assembles the cart view (line totals + subtotal) from the
// joined cart item rows. Money math is kept in int64 cents.
func (s *svc) buildCartResponse(ctx context.Context, cartID int64) (CartResponse, error) {
	items, err := s.repo.ListCartItemsWithVariant(ctx, cartID)
	if err != nil {
		return CartResponse{}, err
	}

	lines := make([]CartLineItem, 0, len(items))
	var subtotal int64

	for _, item := range items {
		unitPrice := int64(item.PriceInCents)
		lineTotal := unitPrice * int64(item.Quantity)
		subtotal += lineTotal

		lines = append(lines, CartLineItem{
			VariantID:        item.VariantID,
			SKU:              item.Sku,
			ProductName:      item.ProductName,
			Quantity:         item.Quantity,
			UnitPriceInCents: unitPrice,
			LineTotalInCents: lineTotal,
		})
	}

	return CartResponse{
		Items:           lines,
		SubtotalInCents: subtotal,
		ItemCount:       len(lines),
	}, nil
}
