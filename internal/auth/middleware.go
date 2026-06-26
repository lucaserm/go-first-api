package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/ecom/internal/json"
	"github.com/lucaserm/ecom/internal/utils"
)

type contextKey string

const userIDContextKey contextKey = "userID"

var (
	errUnauthorized = errors.New("unauthorized")
	errForbidden    = errors.New("admin access required")
)

// Middleware authenticates requests using a Bearer JWT and stores the
// authenticated user id in the request context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		token, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || strings.TrimSpace(token) == "" {
			json.WriteError(w, http.StatusUnauthorized, errUnauthorized)
			return
		}

		userID, err := utils.VerifyJWT(strings.TrimSpace(token))
		if err != nil {
			json.WriteError(w, http.StatusUnauthorized, errUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext returns the authenticated user id stored by Middleware.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

// RequireAdmin must be used after Middleware. It loads the authenticated user
// and rejects the request unless their role is "admin".
func RequireAdmin(queries *repo.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				json.WriteError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			id, err := uuid.Parse(userID)
			if err != nil {
				json.WriteError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			user, err := queries.GetUserByID(r.Context(), pgtype.UUID{Bytes: id, Valid: true})
			if err != nil {
				json.WriteError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			if user.Role != "admin" {
				json.WriteError(w, http.StatusForbidden, errForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
