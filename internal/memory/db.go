package memory

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() error {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getEnv("PG_USER", "asai_user"),
		getEnv("PG_PASS", "secretpassword"),
		getEnv("PG_HOST", "localhost"),
		getEnv("PG_PORT", "5432"),
		getEnv("PG_DB", "asai_db"),
	)

	var err error
	DB, err = pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
