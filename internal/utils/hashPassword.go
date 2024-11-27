package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	const op = "utils.HashPassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", fmt.Errorf("op: %s, err: %w", op, err)
	}
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}