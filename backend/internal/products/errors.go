package products

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product does not have enough stock")
)
