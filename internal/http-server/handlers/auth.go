package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
	dberrors "todoApi/internal/storage/dbErrors"
	"todoApi/internal/utils"
	jwtservice "todoApi/internal/utils/jwt"
)

type JwtResponse struct {
	Token string `json:"jwt"`
}

func Registration(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Login"
		log := log.With(slog.String("op", op))

		requestBody, err := utils.ParseAuthBody(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing body", logger.Err(err))
			return
		}
		if len(requestBody.Password) < 8 || len(requestBody.Username) < 4 {
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
		if err != nil {
			if errors.Is(err, dberrors.ErrAlreadyExists) {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, "User with such username already exists")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while creating user", logger.Err(err))
			return
		}

		token, err := jwtservice.GenerateJwt(createdUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while generating jwt", logger.Err(err))
			return
		}

		jsonResponse, err := json.Marshal(JwtResponse{Token: token})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while marshaling json", logger.Err(err))
			return
		}
		w.Write(jsonResponse)
	})
}

func Login(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func Auth(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
