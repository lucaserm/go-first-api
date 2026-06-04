package orders

import "errors"

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrCustomerIdIsRequired = errors.New("customer ID is required")
	ErrItemsIsRequired      = errors.New("at least one item is required")
)
