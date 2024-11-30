package utils

import (
	"database/sql"
	"fmt"
	"todoApi/internal/storage/models"
)

func ScanUser(rows *sql.Rows) (*models.User, error) {
	const op = "utils.ScanUser"

	var user models.User
	if err := rows.Scan(&user.Id, &user.Username); err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return &user, nil
}
