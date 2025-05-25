package memory

import (
	"asai/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func Init() error {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.AppConfig.Postgres.User,
		config.AppConfig.Postgres.Pass,
		config.AppConfig.Postgres.Host,
		config.AppConfig.Postgres.Port,
		config.AppConfig.Postgres.DB,
	)

	var err error
	ctx := context.Background()
	DB, err = pgxpool.New(ctx, url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := ensureSchema(ctx, DB); err != nil {
		log.Fatalf("Failed to ensure schema: %v", err)
	}

	return nil
}

func ensureSchema(ctx context.Context, db *pgxpool.Pool) error {
	schema := `
CREATE TABLE IF NOT EXISTS embeddings (
	id SERIAL PRIMARY KEY,
	user_id TEXT NOT NULL,
	content TEXT NOT NULL,
	embedding VECTOR(1536),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

	_, err := db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector`)
	if err != nil {
		return fmt.Errorf("failed to enable pgvector extension: %w", err)
	}

	_, err = db.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to ensure schema: %w", err)
	}

	log.Println("Schema ensured (embeddings table exists or created)")
	return nil
}
