package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
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
