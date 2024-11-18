package models

import "time"

type Task struct {
	Id          int       `json:"id"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"is_completed"`
	Deadline    time.Time `json:"deadline_date"`
}
