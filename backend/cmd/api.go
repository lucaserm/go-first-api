package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/ecom/internal/addresses"
	"github.com/lucaserm/ecom/internal/auth"
	"github.com/lucaserm/ecom/internal/cart"
	"github.com/lucaserm/ecom/internal/json"
	"github.com/lucaserm/ecom/internal/orders"
	"github.com/lucaserm/ecom/internal/products"
)

type application struct {
	config config
	db     *pgxpool.Pool
	// logger
}

type config struct {
	addr                string
	db                  dbConfig
	jwtSecret           string
	stripeSecretKey     string
	stripeWebhookSecret string
	easypostAPIKey      string
	corsOrigins         []string
}

type dbConfig struct {
	dsn string
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)                              // rate limiting
	r.Use(middleware.ClientIPFromHeader("CF-Connection-IP")) // rate limiting, analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.config.corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.Write(w, 200, map[string]string{
			"status":  "ok",
			"service": "ecom",
		})
	})

	// Public routes
	productService := products.NewService(repo.New(app.db), app.db)
	productHandler := products.NewHandler(productService)
	productHandler.RegisterRoutes(r)

	authService := auth.NewService(repo.New(app.db))
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(r)

	// Authenticated routes
	orderService := orders.NewService(repo.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService)

	addressService := addresses.NewService(repo.New(app.db), app.db)
	addressHandler := addresses.NewHandler(addressService)

	cartService := cart.NewService(repo.New(app.db), app.db)
	cartHandler := cart.NewHandler(cartService)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware)
		orderHandler.RegisterRoutes(r)
		addressHandler.RegisterRoutes(r)
		cartHandler.RegisterRoutes(r)

		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAdmin(repo.New(app.db)))
			productHandler.RegisterProtectedRoutes(r)
		})
	})

	return r
}

// run -> graceful shutdown
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
