package dberrors

import "errors"

var (
	ErrAlreadyExists = errors.New("record already exists")
	ErrNotFound = errors.New("record not found")
)