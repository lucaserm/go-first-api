package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/lucaserm/ecom/internal/json"
	"github.com/lucaserm/ecom/internal/utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(router *chi.Mux) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", errors))
		return
	}

	user, err := h.service.register(r.Context(), payload)

	if err != nil {
		log.Println(err)
		if err == ErrEmailConflict {
			json.WriteError(w, http.StatusConflict, err)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusCreated, user)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", errors))
		return
	}

	user, err := h.service.login(r.Context(), payload)

	if err != nil {
		log.Println(err)
		if err == ErrInvalidCredentials {
			json.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := utils.CreateJWT(user.ID.String())
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"user":        user,
		"accessToken": token,
	})
}
