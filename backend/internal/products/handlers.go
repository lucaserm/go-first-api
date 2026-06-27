package products

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/lucaserm/ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// RegisterRoutes registers the public, unauthenticated catalog reads.
func (h *handler) RegisterRoutes(router chi.Router) {
	router.Get("/products", h.listProducts)
	router.Get("/products/{id}", h.getProduct)
	router.Get("/categories", h.listCategories)
}

// RegisterProtectedRoutes registers the admin catalog writes. These are mounted
// behind the auth middleware in cmd/api.go.
func (h *handler) RegisterProtectedRoutes(router chi.Router) {
	router.Post("/products", h.createProduct)
	router.Post("/products/{id}/variants", h.createVariant)
	router.Post("/categories", h.createCategory)
}

func (h *handler) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListActiveProducts(r.Context())
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"products": products,
	})
}

func (h *handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.service.GetProductDetail(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"categories": categories,
	})
}

func (h *handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var payload CreateProductPayload
	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", validationErrors))
		return
	}

	product, err := h.service.CreateProduct(r.Context(), payload)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusCreated, product)
}

func (h *handler) createVariant(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload CreateVariantPayload
	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", validationErrors))
		return
	}

	variant, err := h.service.CreateVariant(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusCreated, variant)
}

func (h *handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var payload CreateCategoryPayload
	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", validationErrors))
		return
	}

	category, err := h.service.CreateCategory(r.Context(), payload)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusCreated, category)
}

func parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid product id")
	}
	return id, nil
}
