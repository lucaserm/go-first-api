package products

import (
	"errors"
	"log"
	"net/http"

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
	id := json.Read(r)
	product, err := h.service.GetProductById(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
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
