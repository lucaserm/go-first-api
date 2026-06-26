package products

import (
	"context"

	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
)

// CreateProductPayload is the admin request body for creating a product.
type CreateProductPayload struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=draft active archived"`
	CategoryID  *int64 `json:"categoryId"`
}

// CreateVariantPayload is the admin request body for creating a SKU/variant
// together with the option values it is composed of.
type CreateVariantPayload struct {
	SKU            string  `json:"sku" validate:"required"`
	PriceInCents   int32   `json:"priceInCents" validate:"gte=0"`
	Stock          int32   `json:"stock" validate:"gte=0"`
	WeightGrams    int32   `json:"weightGrams" validate:"gte=0"`
	OptionValueIDs []int64 `json:"optionValueIds"`
}

// CreateCategoryPayload is the admin request body for creating a category.
type CreateCategoryPayload struct {
	Name     string `json:"name" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	ParentID *int64 `json:"parentId"`
}

// ProductDetail is the public read aggregate for a single product: the product
// itself plus its variants, options and images.
type ProductDetail struct {
	Product  repo.Product          `json:"product"`
	Variants []repo.ProductVariant `json:"variants"`
	Options  []repo.ProductOption  `json:"options"`
	Images   []repo.ProductImage   `json:"images"`
}

type Service interface {
	ListActiveProducts(ctx context.Context) ([]repo.Product, error)
	GetProductDetail(ctx context.Context, id int64) (ProductDetail, error)
	CreateProduct(ctx context.Context, payload CreateProductPayload) (repo.Product, error)
	CreateVariant(ctx context.Context, productID int64, payload CreateVariantPayload) (repo.ProductVariant, error)
	ListCategories(ctx context.Context) ([]repo.Category, error)
	CreateCategory(ctx context.Context, payload CreateCategoryPayload) (repo.Category, error)
}
