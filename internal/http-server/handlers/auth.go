package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
	dberrors "todoApi/internal/storage/dbErrors"
	"todoApi/internal/utils"
)

func Registration(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Login"
		log := log.With(slog.String("op", op))

		requestBody, err := utils.ParseAuthBody(r)
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
		
		passwordHash, err := utils.HashPassword(requestBody.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while hashing password", logger.Err(err))
			return
		}
		
		createdUserID, err := s.CreateUser(requestBody.Username, passwordHash)
		if err != nil{
			if errors.Is(err, dberrors.ErrAlreadyExists){
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while creating user", logger.Err(err))
			return
		}
		_ = createdUserID
		//TODO: return jwt token
	})
}

func Login(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}

func Auth(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}