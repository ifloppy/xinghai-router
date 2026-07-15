package app

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Service struct {
	cfg        Config
	db         *pgxpool.Pool
	httpClient *http.Client
	limiter    *limiter
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	if err := migrate(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return &Service{cfg: cfg, db: db, httpClient: &http.Client{Timeout: cfg.RequestTimeout}, limiter: newLimiter(cfg.RateLimitPerMinute)}, nil
}
func (s *Service) Close()                { s.db.Close() }
func (s *Service) Handler() http.Handler { return s.routes() }
