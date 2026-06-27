package cart

import "errors"

var (
	ErrVariantNotFound   = errors.New("variant not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrCartItemNotFound  = errors.New("cart item not found")
)
