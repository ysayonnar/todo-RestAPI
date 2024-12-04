package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
	dberrors "todoApi/internal/storage/dbErrors"
	"todoApi/internal/storage/models"
	"todoApi/internal/utils"
)

type AllUsersResponse struct {
	Users []models.User `json:"users"`
}

func GetAllUsers(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetAllUsers"
		log := log.With(slog.String("op", op))

		rows, err := s.GetAllUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while database request", logger.Err(err))
			return
		}

		var users AllUsersResponse
		for rows.Next() {
			user, err := utils.ScanUser(rows)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("Error while scaning users", logger.Err(err))
				return
			}
			users.Users = append(users.Users, *user)
		}
		defer rows.Close()

		jsonResponse, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing json", logger.Err(err))
			return
		}

		w.Write(jsonResponse)
	})
}

func GetUserById(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetUserById"
		log := log.With(slog.String("op", op))

		stringId := r.Header.Get("user-id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing id from header to int", logger.Err(err))
			return
		}

		rows, err := s.GetUserById(id)
		if err != nil {
			if errors.Is(err, dberrors.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while db request", logger.Err(err))
			return
		}
		defer rows.Close()

		isExists := rows.Next()
		if !isExists {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		user, err := utils.ScanUser(rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while scanning rows", logger.Err(err))
			return
		}

		jsonResponse, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing json", logger.Err(err))
			return
		}

		w.Write(jsonResponse)
	})
}
