package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	jwtservice "todoApi/internal/utils/jwt"

	"github.com/golang-jwt/jwt/v5"
)

func AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "jwt token required, login first")
			return
		}
		jwtToken, err := jwtservice.ValidateJWT(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !jwtToken.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, ok := claims["user_id"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId = int(userId.(float64)) //type assertion for convert
		//TODO: переписать через контекст

		r.Header.Set("user-id", strconv.Itoa(userId.(int)))
		next.ServeHTTP(w, r)
	})
}
