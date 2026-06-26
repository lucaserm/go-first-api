package orders

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lucaserm/ecom/internal/auth"
	"github.com/lucaserm/ecom/internal/json"
	"github.com/lucaserm/ecom/internal/products"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(router chi.Router) {
	router.Post("/orders", h.PlaceOrder)
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.WriteError(w, http.StatusUnauthorized, ErrCustomerIdIsRequired)
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		json.WriteError(w, http.StatusUnauthorized, ErrCustomerIdIsRequired)
		return
	}

	customerID := pgtype.UUID{Bytes: parsedID, Valid: true}

	var payload createOrderParams

	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), customerID, payload)

	if err != nil {
		log.Println(err)

		if errors.Is(err, products.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
