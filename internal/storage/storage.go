package storage

import (
	"database/sql"
	"fmt"
	"time"
	"todoApi/internal/config"
	"todoApi/internal/storage/queries"

	_ "github.com/lib/pq"
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

func (s *Storage) CreateTask(taskText string, deadline time.Time) (int, error) {
	const op = "storage.CreateTask"

	var id int
	query := `INSERT INTO tasks (task, deadline_date) VALUES ($1, $2) RETURNING id;`
	err := s.db.QueryRow(query, taskText, deadline).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetAllTasks() (*sql.Rows, error) {
	const op = "storage.GetAllTasks"

	query := `SELECT id, task, is_completed, deadline_date FROM tasks ORDER BY deadline_date`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}
	return rows, nil
}

func (s *Storage) GetTaskById(id int) (*sql.Rows, error) {
	const op = "storage.GetTaskById"

	query := `SELECT id, task, is_completed, deadline_date FROM tasks WHERE id = $1;`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}
	return rows, nil
}

func (s *Storage) DeleteTaskById(id int) error {
	const op = "storage.DeleteTaskById"

	query := `DELETE FROM tasks WHERE id = $1;`

	rows, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}
	if affected < 1 {
		return fmt.Errorf("op: %s, err: not found", op)
	}
	return nil
}

func (s *Storage) SetTaskCompletedById(id int) error {
	const op = "storage.SetTaskCompletedById"

	query := `UPDATE tasks SET is_completed = true WHERE id = $1;`
	rows, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}
	if affected < 1 {
		return fmt.Errorf("op: %s, err: not found", op)
	}
	return nil
}

func (s *Storage) GetUncomplitedTasks() (*sql.Rows, error) {
	const op = "storage.GetTomorowTasks"

	query := `SELECT id, task, is_completed, deadline_date FROM tasks WHERE deadline_date < $1 AND is_completed = false;`
	rows, err := s.db.Query(query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return rows, nil
}

func(s *Storage) GetTodaysTasks() (*sql.Rows, error){
	const op = "storage.GetTodaysTasks"

	query := `SELECT id, task, is_completed, deadline_date FROM tasks WHERE deadline_date = $1 AND is_completed = false;`

	rows,err := s.db.Query(query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}
	return rows, nil
}
