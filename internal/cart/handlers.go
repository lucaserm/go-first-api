package cart

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
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// RegisterRoutes registers the authenticated, user-scoped cart routes. These
// are mounted behind the auth middleware in cmd/api.go.
func (h *handler) RegisterRoutes(router chi.Router) {
	router.Get("/cart", h.getCart)
	router.Post("/cart/items", h.addItem)
	router.Patch("/cart/items/{variantId}", h.updateItem)
	router.Delete("/cart/items/{variantId}", h.removeItem)
	router.Delete("/cart", h.clearCart)
}

func (h *handler) getCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	cart, err := h.service.GetCart(r.Context(), userID)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, cart)
}

func (h *handler) addItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	var payload AddItemPayload
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

	cart, err := h.service.AddItem(r.Context(), userID, payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.Write(w, http.StatusCreated, cart)
}

func (h *handler) updateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	variantID, err := parseVariantID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload UpdateItemPayload
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

	cart, err := h.service.UpdateItem(r.Context(), userID, variantID, payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.Write(w, http.StatusOK, cart)
}

func (h *handler) removeItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	variantID, err := parseVariantID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.RemoveItem(r.Context(), userID, variantID); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) clearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	if err := h.service.ClearCart(r.Context(), userID); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeServiceError maps cart service sentinel errors to HTTP statuses.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrVariantNotFound), errors.Is(err, ErrCartItemNotFound):
		json.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, ErrInsufficientStock):
		json.WriteError(w, http.StatusConflict, err)
	default:
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
	}
}

// userIDFromRequest resolves the authenticated user id from the request context
// and parses it into a pgtype.UUID. It writes a 401 response and returns false
// when the user is not authenticated or the id is malformed.
func userIDFromRequest(w http.ResponseWriter, r *http.Request) (pgtype.UUID, bool) {
	rawID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return pgtype.UUID{}, false
	}

	parsedID, err := uuid.Parse(rawID)
	if err != nil {
		json.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return pgtype.UUID{}, false
	}

	return pgtype.UUID{Bytes: parsedID, Valid: true}, true
}

func parseVariantID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "variantId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid variant id")
	}
	return id, nil
}
