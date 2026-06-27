package addresses

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

// RegisterRoutes registers the authenticated, user-scoped address routes. These
// are mounted behind the auth middleware in cmd/api.go.
func (h *handler) RegisterRoutes(router chi.Router) {
	router.Get("/addresses", h.listAddresses)
	router.Post("/addresses", h.createAddress)
	router.Patch("/addresses/{id}", h.updateAddress)
	router.Delete("/addresses/{id}", h.deleteAddress)
}

func (h *handler) listAddresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"addresses": addresses,
	})
}

func (h *handler) createAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	var payload CreateAddressPayload
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

	address, err := h.service.CreateAddress(r.Context(), userID, payload)
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusCreated, address)
}

func (h *handler) updateAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, err := parseID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload UpdateAddressPayload
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

	address, err := h.service.UpdateAddress(r.Context(), userID, id, payload)
	if err != nil {
		if errors.Is(err, ErrAddressNotFound) {
			json.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, address)
}

func (h *handler) deleteAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, err := parseID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.DeleteAddress(r.Context(), userID, id); err != nil {
		if errors.Is(err, ErrAddressNotFound) {
			json.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid address id")
	}
	return id, nil
}
