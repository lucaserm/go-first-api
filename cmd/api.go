package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/lucaserm/ecom/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/ecom/internal/products"
)

type application struct {
	config config
	db     *pgx.Conn
	// logger
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)                              // rate limiting
	r.Use(middleware.ClientIPFromHeader("CF-Connection-IP")) // rate limiting, analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok."))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)

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

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
