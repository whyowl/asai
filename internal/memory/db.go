package memory

import (
	"asai/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func Init(ctx context.Context, dimension int) error {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.AppConfig.Postgres.User,
		config.AppConfig.Postgres.Pass,
		config.AppConfig.Postgres.Host,
		config.AppConfig.Postgres.Port,
		config.AppConfig.Postgres.DB,
	)

	var err error
	DB, err = pgxpool.New(ctx, url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := ensureSchema(ctx, DB, "memory", dimension); err != nil {
		log.Fatalf("Failed to ensure schema: %v", err)
	}

	return nil
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool, tableName string, vectorDim int) error {
	var exists bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = $1
		)
	`, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if exists {
		var actualDimStr string
		err = pool.QueryRow(ctx, fmt.Sprintf(`
			SELECT pg_catalog.format_type(a.atttypid, a.atttypmod)
			FROM pg_catalog.pg_attribute a
			JOIN pg_catalog.pg_class c ON a.attrelid = c.oid
			WHERE c.relname = $1 AND a.attname = 'embedding' AND a.attnum > 0 AND NOT a.attisdropped
		`), tableName).Scan(&actualDimStr)
		if err != nil {
			return fmt.Errorf("failed to read embedding column type: %w", err)
		}

		var actualDim int
		_, err := fmt.Sscanf(actualDimStr, "vector(%d)", &actualDim)
		if err != nil {
			return fmt.Errorf("unexpected format of embedding column: %v", actualDimStr)
		}

		if actualDim != vectorDim {
			return fmt.Errorf("existing embedding dimension (%d) != expected (%d), you need drop table manually", actualDim, vectorDim)
		}

		log.Println("Table exists with correct embedding dimension.")
		return nil
	}

	createTableSQL := fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			user_id TEXT,
			text TEXT,
			embedding VECTOR(%d)
		);
	`, tableName, vectorDim)

	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Printf("Created table %s with embedding dimension %d.\n", tableName, vectorDim)
	return nil
}
