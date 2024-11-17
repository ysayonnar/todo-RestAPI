package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
)

type TaskCreateRequest struct {
	Task         string `json:"task"`
	DeadlineDate string `json:"deadline"`
}

type TaskCreateResponse struct {
	Id int `json:"id"`
}

func CreateTask(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.CreateTask"
		log := log.With(
			slog.String("op", op),
		)

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error while parsing body", logger.Err(err))
			fmt.Fprint(w, "error!")
		}
		defer r.Body.Close()

		var request TaskCreateRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("error while parsing json")
			fmt.Fprint(w, "invalid body")
			return
		}
		if len(request.DeadlineDate) == 0 || len(request.Task) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid body")
			return
		}

		//01/02/2006 как я понял это дефолтsный паттерн
		deadlineTime, err := time.Parse("01/02/2006", request.DeadlineDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid date")
			return
		}

		id, err := s.CreateTask(request.Task, deadlineTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("err", logger.Err(err))
			return
		}

		jsonResponse, err := json.Marshal(&TaskCreateResponse{Id: id})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error while parsing json response", logger.Err(err))
			return
		}
		w.Write(jsonResponse)
	})
}
