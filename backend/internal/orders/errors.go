package orders

import "errors"

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrCustomerIdIsRequired    = errors.New("customer ID is required")
	ErrAddressNotFound         = errors.New("address not found")
	ErrCartEmpty               = errors.New("cart is empty")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
)
