package response

import (
	"database/sql"
	"fmt"
	"todoApi/internal/storage/models"
)

func ScanTask(rows *sql.Rows) (*models.Task, error){
	const op = "response.ScanTask"

	var task models.Task
	if err := rows.Scan(&task.Id, &task.Task, &task.IsCompleted, &task.Deadline); err != nil{
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return &task, nil
}