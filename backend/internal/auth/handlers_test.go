package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAuthServiceHandlers(t *testing.T) {
	service := &mockAuthService{}
	handler := NewHandler(service)

	t.Run("should fail if the user payload is null", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/register", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := RegisterPayload{
			Username: "john",
			Email:    "invalid",
			Password: "password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := RegisterPayload{
			Username: "john",
			Email:    "valid@mail.com",
			Password: "password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should correctly login the user", func(t *testing.T) {
		payload := LoginPayload{
			Email:    "valid@mail.com",
			Password: "password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockAuthService struct{}

func (m *mockAuthService) register(ctx context.Context, payload RegisterPayload) (UserResponse, error) {
	return UserResponse{}, nil
}
func (m *mockAuthService) login(ctx context.Context, payload LoginPayload) (UserResponse, error) {
	return UserResponse{}, nil
}
