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

// PlaceOrder builds an order from the authenticated user's cart, snapshotting
// the chosen shipping address and each line item so the order is immutable
// against later catalog/address changes.
//
// NOTE: stock is decremented here, at creation, FOR NOW because payment is not
// yet built. In Phase 5 the stock decrement moves to the Stripe payment-success
// webhook and the creation status becomes 'awaiting_payment' instead of 'pending'.
// Shipping and tax totals are 0 placeholders until Phases 5/6.
func (s *svc) PlaceOrder(ctx context.Context, customerID pgtype.UUID, payload CreateOrderPayload) (OrderResponse, error) {
	if !customerID.Valid {
		return OrderResponse{}, ErrCustomerIdIsRequired
	}

	// 1. Resolve and snapshot the shipping address (must belong to this user).
	address, err := s.repo.GetAddressByIDForUser(ctx, repo.GetAddressByIDForUserParams{
		ID:     payload.AddressID,
		UserID: customerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderResponse{}, ErrAddressNotFound
		}
		return OrderResponse{}, err
	}

	// 2. Load the user's cart lines.
	cart, err := s.repo.GetCartByUser(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderResponse{}, ErrCartEmpty
		}
		return OrderResponse{}, err
	}

	cartItems, err := s.repo.ListCartItemsWithVariant(ctx, cart.ID)
	if err != nil {
		return OrderResponse{}, err
	}
	if len(cartItems) == 0 {
		return OrderResponse{}, ErrCartEmpty
	}

	// Compute subtotal and re-validate stock before opening the transaction.
	var subtotal int64
	for _, line := range cartItems {
		if line.Quantity > line.Stock {
			return OrderResponse{}, products.ErrProductNoStock
		}
		subtotal += int64(line.PriceInCents) * int64(line.Quantity)
	}

	// 3. Persist atomically: create order, snapshot items, decrement stock.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return OrderResponse{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, repo.CreateOrderParams{
		CustomerID:        customerID,
		Status:            StatusPending,
		Currency:          "usd",
		SubtotalCents:     int32(subtotal),
		ShippingCents:     0, // placeholder until Phase 6 (shipping)
		TaxCents:          0, // placeholder until Phase 6 (tax)
		TotalCents:        int32(subtotal),
		ShippingAddressID: pgtype.Int8{Int64: address.ID, Valid: true},
		ShipRecipientName: address.RecipientName,
		ShipLine1:         address.Line1,
		ShipLine2:         address.Line2,
		ShipCity:          address.City,
		ShipRegion:        address.Region,
		ShipPostalCode:    address.PostalCode,
		ShipCountry:       address.Country,
		ShipPhone:         address.Phone,
	})
	if err != nil {
		return OrderResponse{}, err
	}

	items := make([]OrderLineItem, 0, len(cartItems))
	for _, line := range cartItems {
		if _, err := qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      order.ID,
			VariantID:    line.VariantID,
			Quantity:     line.Quantity,
			PriceInCents: line.PriceInCents,
			VariantSku:   line.Sku,
			ProductName:  line.ProductName,
		}); err != nil {
			return OrderResponse{}, err
		}

		// Atomic stock guard: DecreaseVariantStock only updates when stock >= qty.
		// Zero rows affected means another writer drained the stock concurrently,
		// so abort the whole order.
		rows, err := qtx.DecreaseVariantStock(ctx, repo.DecreaseVariantStockParams{
			Stock: line.Quantity,
			ID:    line.VariantID,
		})
		if err != nil {
			return OrderResponse{}, err
		}
		if rows == 0 {
			return OrderResponse{}, products.ErrProductNoStock
		}

		items = append(items, OrderLineItem{
			VariantID:        line.VariantID,
			SKU:              line.Sku,
			ProductName:      line.ProductName,
			Quantity:         line.Quantity,
			UnitPriceInCents: int64(line.PriceInCents),
			LineTotalInCents: int64(line.PriceInCents) * int64(line.Quantity),
		})
	}

	// 4. Empty the cart now that its contents have become an order.
	if err := qtx.ClearCart(ctx, cart.ID); err != nil {
		return OrderResponse{}, err
	}

	// 5. Commit.
	if err := tx.Commit(ctx); err != nil {
		return OrderResponse{}, err
	}

	return buildOrderResponse(order, items), nil
}

func (s *svc) GetOrder(ctx context.Context, customerID pgtype.UUID, orderID int64) (OrderResponse, error) {
	order, err := s.repo.GetOrderByIDForCustomer(ctx, repo.GetOrderByIDForCustomerParams{
		ID:         orderID,
		CustomerID: customerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderResponse{}, ErrOrderNotFound
		}
		return OrderResponse{}, err
	}

	items, err := s.loadItems(ctx, order.ID)
	if err != nil {
		return OrderResponse{}, err
	}

	return buildOrderResponse(order, items), nil
}

func (s *svc) ListOrders(ctx context.Context, customerID pgtype.UUID) ([]OrderResponse, error) {
	rows, err := s.repo.ListOrdersByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	responses := make([]OrderResponse, 0, len(rows))
	for _, order := range rows {
		items, err := s.loadItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		responses = append(responses, buildOrderResponse(order, items))
	}

	return responses, nil
}

// UpdateStatus validates the requested lifecycle transition for one of the
// customer's orders and persists it. It is callable now (e.g. from an admin
// route) and is also the seam Phase 5 uses to drive status from payment events.
func (s *svc) UpdateStatus(ctx context.Context, customerID pgtype.UUID, orderID int64, newStatus string) (OrderResponse, error) {
	order, err := s.repo.GetOrderByIDForCustomer(ctx, repo.GetOrderByIDForCustomerParams{
		ID:         orderID,
		CustomerID: customerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderResponse{}, ErrOrderNotFound
		}
		return OrderResponse{}, err
	}

	if !canTransition(order.Status, newStatus) {
		return OrderResponse{}, ErrInvalidStatusTransition
	}

	updated, err := s.repo.UpdateOrderStatus(ctx, repo.UpdateOrderStatusParams{
		ID:     order.ID,
		Status: newStatus,
	})
	if err != nil {
		return OrderResponse{}, err
	}

	items, err := s.loadItems(ctx, updated.ID)
	if err != nil {
		return OrderResponse{}, err
	}

	return buildOrderResponse(updated, items), nil
}

func (s *svc) loadItems(ctx context.Context, orderID int64) ([]OrderLineItem, error) {
	rows, err := s.repo.ListOrderItemsByOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items := make([]OrderLineItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, OrderLineItem{
			VariantID:        row.VariantID,
			SKU:              row.VariantSku,
			ProductName:      row.ProductName,
			Quantity:         row.Quantity,
			UnitPriceInCents: int64(row.PriceInCents),
			LineTotalInCents: int64(row.PriceInCents) * int64(row.Quantity),
		})
	}

	return items, nil
}

func buildOrderResponse(order repo.Order, items []OrderLineItem) OrderResponse {
	return OrderResponse{
		ID:            order.ID,
		Status:        order.Status,
		Currency:      order.Currency,
		SubtotalCents: int64(order.SubtotalCents),
		ShippingCents: int64(order.ShippingCents),
		TaxCents:      int64(order.TaxCents),
		TotalCents:    int64(order.TotalCents),
		Shipping: ShippingAddress{
			RecipientName: order.ShipRecipientName,
			Line1:         order.ShipLine1,
			Line2:         order.ShipLine2,
			City:          order.ShipCity,
			Region:        order.ShipRegion,
			PostalCode:    order.ShipPostalCode,
			Country:       order.ShipCountry,
			Phone:         order.ShipPhone,
		},
		Items: items,
	}
}
