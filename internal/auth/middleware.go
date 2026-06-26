package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/lucaserm/ecom/internal/json"
	"github.com/lucaserm/ecom/internal/utils"
)

type contextKey string

const userIDContextKey contextKey = "userID"

var errUnauthorized = errors.New("unauthorized")

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
