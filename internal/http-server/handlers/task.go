package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
	"todoApi/internal/storage/models"

	"github.com/gorilla/mux"
)

type TaskCreateRequest struct {
	Task         string `json:"task"`
	DeadlineDate string `json:"deadline"`
}

type TaskCreateResponse struct {
	Id int `json:"id"`
}

type TaskAllResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func CreateTask(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.CreateTask"
		log = log.With(
			slog.String("op", op),
		)

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

		//02/01/2006 как я понял это дефолтsный паттерн
		deadlineTime, err := time.Parse("02/01/2006", request.DeadlineDate)
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

func GetAllTasks(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetAllTasks"
		log = log.With(slog.String("op", op))

		rows, err := s.GetAllTasks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("db error", logger.Err(err))
			return
		}
		defer rows.Close()

		var tasks TaskAllResponse

		for rows.Next() {
			var id int
			var task string
			var is_completed bool
			var deadline_date time.Time

			if err := rows.Scan(&id, &task, &is_completed, &deadline_date); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("error while scaning rows", logger.Err(err))
				return
			}

			item := models.Task{
				Id:          id,
				Task:        task,
				IsCompleted: is_completed,
				Deadline:    deadline_date,
			}
			tasks.Tasks = append(tasks.Tasks, item)
		}

		jsonResponse, err := json.Marshal(tasks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error while parsing jsonResponse", logger.Err(err))
			return
		}
		w.Write(jsonResponse)
	})
}

func GetTaskById(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTaskById"
		log := log.With(slog.String("op", op))

		params := mux.Vars(r)
		if _, ok := params["id"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid id")
			return
		}

		taskId, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid id")
			return
		}
		rows, err := s.GetTaskById(taskId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while db request", logger.Err(err))
			return
		}
		defer rows.Close()

		var id int
		var task string
		var is_completed bool
		var deadline_date time.Time
		isExists := rows.Next()
		if !isExists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no task with such id")
			return
		}
		if err := rows.Scan(&id, &task, &is_completed, &deadline_date); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error while scaning rows", logger.Err(err))
			return
		}

		item := models.Task{
			Id:          id,
			Task:        task,
			IsCompleted: is_completed,
			Deadline:    deadline_date,
		}

		jsonResponse, err := json.Marshal(item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing jsonResponse", logger.Err(err))
			return
		}
		w.Write(jsonResponse)
	})
}

func DeleteTaskById(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeleteTaskById"

		log := log.With(slog.String("op", op))
		params := mux.Vars(r)
		if _, ok := params["id"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid id")
			return
		}

		taskId, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid id")
			return
		}

		err = s.DeleteTaskById(taskId)
		if err != nil {
			log.Error("error while deleting task", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "deleted")
	})
}

func SetTaskCompletedById(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.SetTaskCompletedById"

		log := log.With(slog.String("op", op))
		params := mux.Vars(r)
		if _, ok := params["id"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid id")
			return
		}

		taskId, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid id")
			return
		}

		err = s.SetTaskCompletedById(taskId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("error while seting is_completed", logger.Err(err))
			return
		}
		fmt.Fprint(w, "is_completed setted")
	})
}

func GetUncomplitedTasks(log *slog.Logger, s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTomorowTasks"

		log := log.With(slog.String("op", op))
		rows, err := s.GetUncomplitedTasks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Database error", logger.Err(err))
			return
		}
		defer rows.Close()

		tasks := TaskAllResponse{}
		for rows.Next() {
			var id int
			var task string
			var is_completed bool
			var deadline_date time.Time

			if err := rows.Scan(&id, &task, &is_completed, &deadline_date); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("error while scaning rows", logger.Err(err))
				return
			}

			item := models.Task{
				Id:          id,
				Task:        task,
				IsCompleted: is_completed,
				Deadline:    deadline_date,
			}
			tasks.Tasks = append(tasks.Tasks, item)
		}

		jsonResponse, err := json.Marshal(tasks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error while parsing jsonResponse", logger.Err(err))
			return
		}
		w.Write(jsonResponse)
	})
}

func GetTodaysTasks(log *slog.Logger, s *storage.Storage) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTodaysTasks"
		log := log.With(slog.String("op", op))

		rows, err := s.GetTodaysTasks()
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while database request", logger.Err(err))
			return
		}


		//TODO: отдельную функцию для парсинга sql.rows + использовать ее во всех роутах где она нужна

	})
}
