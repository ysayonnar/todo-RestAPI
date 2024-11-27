package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
	"todoApi/internal/storage/models"
	"todoApi/internal/utils"

	"github.com/gorilla/mux"
)

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

		request, err := utils.ParseTaskBody(r)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing body", logger.Err(err))
			return
		}

		if len(request.DeadlineDate) == 0 || len(request.Task) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid body")
			fmt.Fprint(w, request)
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
			item, err := utils.ScanTask(rows)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("Error while scaning rows", logger.Err(err))
				return
			}
			tasks.Tasks = append(tasks.Tasks, *item)
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

		isExists := rows.Next()
		if !isExists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no task with such id")
			return
		}

		item, err := utils.ScanTask(rows)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while scaning rows", logger.Err(err))
			return
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
			item, err := utils.ScanTask(rows)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("Error while parsing sql.rows", logger.Err(err))
				return
			}
			tasks.Tasks = append(tasks.Tasks, *item)
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
		defer rows.Close()
		
		tasksResponse := TaskAllResponse{}	
		for rows.Next(){
			item, err := utils.ScanTask(rows)
			if err != nil{
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("Error while scanning rows", logger.Err(err))
				return
			}
			tasksResponse.Tasks = append(tasksResponse.Tasks, *item)
		}
		jsonResponse, err := json.Marshal(tasksResponse)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error while parsing to json", logger.Err(err))
			return
		}
		w.Write(jsonResponse)
	})
}
