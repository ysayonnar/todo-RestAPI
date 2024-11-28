package jwtservice

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("secret")

func GenerateJwt(userId int) (string, error) {
	const op = "jwtservie.GenerateJwt"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	strToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("op: %s, err: %w", op, err)
	}
	return strToken, nil
}
