package products

import (
	"context"

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

func (s *svc) ListActiveProducts(ctx context.Context) ([]repo.Product, error) {
	products, err := s.repo.ListActiveProducts(ctx)
	if err != nil {
		return []repo.Product{}, err
	}

	if products == nil {
		return []repo.Product{}, nil
	}

	return products, nil
}

func (s *svc) GetProductDetail(ctx context.Context, id int64) (ProductDetail, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ProductDetail{}, ErrProductNotFound
		}
		return ProductDetail{}, err
	}

	variants, err := s.repo.ListVariantsByProduct(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}

	options, err := s.repo.ListOptionsByProduct(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}

	images, err := s.repo.ListImagesByProduct(ctx, id)
	if err != nil {
		return ProductDetail{}, err
	}

	if variants == nil {
		variants = []repo.ProductVariant{}
	}
	if options == nil {
		options = []repo.ProductOption{}
	}
	if images == nil {
		images = []repo.ProductImage{}
	}

	return ProductDetail{
		Product:  product,
		Variants: variants,
		Options:  options,
		Images:   images,
	}, nil
}

func (s *svc) CreateProduct(ctx context.Context, payload CreateProductPayload) (repo.Product, error) {
	status := payload.Status
	if status == "" {
		status = "draft"
	}

	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:        payload.Name,
		Slug:        pgtype.Text{String: payload.Slug, Valid: true},
		Description: payload.Description,
		Status:      status,
		CategoryID:  toPgInt8(payload.CategoryID),
	})
}

func (s *svc) CreateVariant(ctx context.Context, productID int64, payload CreateVariantPayload) (repo.ProductVariant, error) {
	if _, err := s.repo.GetProductByID(ctx, productID); err != nil {
		if err == pgx.ErrNoRows {
			return repo.ProductVariant{}, ErrProductNotFound
		}
		return repo.ProductVariant{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.ProductVariant{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.repo.WithTx(tx)

	variant, err := qtx.CreateProductVariant(ctx, repo.CreateProductVariantParams{
		ProductID:    productID,
		Sku:          payload.SKU,
		PriceInCents: payload.PriceInCents,
		Stock:        payload.Stock,
		WeightGrams:  payload.WeightGrams,
	})
	if err != nil {
		return repo.ProductVariant{}, err
	}

	for _, optionValueID := range payload.OptionValueIDs {
		if err := qtx.LinkVariantOptionValue(ctx, repo.LinkVariantOptionValueParams{
			VariantID:     variant.ID,
			OptionValueID: optionValueID,
		}); err != nil {
			return repo.ProductVariant{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.ProductVariant{}, err
	}

	return variant, nil
}

func (s *svc) ListCategories(ctx context.Context) ([]repo.Category, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return []repo.Category{}, err
	}

	if categories == nil {
		return []repo.Category{}, nil
	}

	return categories, nil
}

func (s *svc) CreateCategory(ctx context.Context, payload CreateCategoryPayload) (repo.Category, error) {
	return s.repo.CreateCategory(ctx, repo.CreateCategoryParams{
		Name:     payload.Name,
		Slug:     payload.Slug,
		ParentID: toPgInt8(payload.ParentID),
	})
}

func toPgInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}
