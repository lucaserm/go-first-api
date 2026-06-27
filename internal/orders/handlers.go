package orders

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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

// RegisterRoutes registers the authenticated, customer-scoped order routes.
// These are mounted behind the auth middleware in cmd/api.go.
func (h *handler) RegisterRoutes(router chi.Router) {
	router.Post("/orders", h.PlaceOrder)
	router.Get("/orders", h.ListOrders)
	router.Get("/orders/{id}", h.GetOrder)
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	customerID, ok := customerIDFromRequest(w, r)
	if !ok {
		return
	}

	var payload CreateOrderPayload
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

	order, err := h.service.PlaceOrder(r.Context(), customerID, payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.Write(w, http.StatusCreated, order)
}

func (h *handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	customerID, ok := customerIDFromRequest(w, r)
	if !ok {
		return
	}

	orders, err := h.service.ListOrders(r.Context(), customerID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.Write(w, http.StatusOK, orders)
}

func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	customerID, ok := customerIDFromRequest(w, r)
	if !ok {
		return
	}

	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order id"))
		return
	}

	order, err := h.service.GetOrder(r.Context(), customerID, orderID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.Write(w, http.StatusOK, order)
}

// writeServiceError maps order service sentinel errors to HTTP statuses.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrAddressNotFound), errors.Is(err, ErrOrderNotFound):
		json.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, ErrCartEmpty):
		json.WriteError(w, http.StatusConflict, err)
	case errors.Is(err, products.ErrProductNoStock):
		json.WriteError(w, http.StatusConflict, err)
	case errors.Is(err, ErrInvalidStatusTransition):
		json.WriteError(w, http.StatusConflict, err)
	default:
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
	}
}

// customerIDFromRequest resolves the authenticated user id from the request
// context and parses it into a pgtype.UUID. It writes a 401 response and
// returns false when the user is not authenticated or the id is malformed.
func customerIDFromRequest(w http.ResponseWriter, r *http.Request) (pgtype.UUID, bool) {
	rawID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.WriteError(w, http.StatusUnauthorized, ErrCustomerIdIsRequired)
		return pgtype.UUID{}, false
	}

	parsedID, err := uuid.Parse(rawID)
	if err != nil {
		json.WriteError(w, http.StatusUnauthorized, ErrCustomerIdIsRequired)
		return pgtype.UUID{}, false
	}

	return pgtype.UUID{Bytes: parsedID, Valid: true}, true
}
