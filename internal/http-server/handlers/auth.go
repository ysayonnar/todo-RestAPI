package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
	"todoApi/internal/utils"
)

func Login(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Login"
		log := log.With(slog.String("op", op))

		requestBody, err := utils.ParseAuthBody(r)
		fmt.Println(requestBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing body", logger.Err(err))
			return
		}

		if len(requestBody.Password) < 8 || len(requestBody.Username) < 4{
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid password or username")
			return
		}

		fmt.Fprint(w, "u all set")
	})
}

func Registration(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}

func Auth(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}