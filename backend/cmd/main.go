package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucaserm/ecom/internal/env"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
		jwtSecret:           env.GetString("JWT_SECRET", ""),
		stripeSecretKey:     env.GetString("STRIPE_SECRET_KEY", ""),
		stripeWebhookSecret: env.GetString("STRIPE_WEBHOOK_SECRET", ""),
		easypostAPIKey:      env.GetString("EASYPOST_API_KEY", ""),
		corsOrigins:         strings.Split(env.GetString("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
	}

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if cfg.jwtSecret == "" {
		slog.Error("JWT_SECRET environment variable is required")
		os.Exit(1)
	}

	// database
	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	logger.Info("connected to database")

	api := application{
		config: cfg,
		db:     pool,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
