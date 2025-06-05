package memory

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

func InsertEmbedding(ctx context.Context, db *pgxpool.Pool, table, userID, text string, embedding []float32) error {
	vec := pgvector.NewVector(embedding)

	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, text, embedding)
		VALUES ($1, $2, $3)
	`, table)

	_, err := db.Exec(ctx, query, userID, text, vec)
	return err
}

func QuerySimilarEmbeddings(ctx context.Context, db *pgxpool.Pool, table, userID string, embedding []float32, topK int) ([]string, error) {
	vec := pgvector.NewVector(embedding)

	query := fmt.Sprintf(`
		SELECT text
		FROM %s
		WHERE user_id = $1
		ORDER BY embedding <-> $2
		LIMIT $3
	`, table)

	rows, err := db.Query(ctx, query, userID, vec, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return nil, err
		}
		results = append(results, text)
	}

	return results, nil
}
