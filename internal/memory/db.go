package memory

import (
	"asai/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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
	DB, err = pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	return nil
}
