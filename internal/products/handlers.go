package products

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
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

func (h *handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductById(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			http.Error(w, "product not found for id "+idStr, http.StatusNotFound)
			return
		}

		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, map[string][]repo.Product{
		"products": products,
	})
}
