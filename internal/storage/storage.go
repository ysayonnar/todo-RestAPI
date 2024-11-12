package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"todoApi/internal/config"
	"todoApi/internal/storage/queries"
)

type Storage struct {
	db *sql.DB
}

func createConnectionString(cfg *config.Postgres) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", cfg.User, cfg.Password, cfg.DBName, cfg.SslMode)
}

func NewStorageConnection(postgresConfig *config.Postgres) (*Storage, error) {
	const op = "storage.NewStorageConnection"

	connectionString := createConnectionString(postgresConfig)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(queries.UsersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(queries.TasksTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
