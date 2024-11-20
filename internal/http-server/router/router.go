package router

import (
	"log/slog"
	"todoApi/internal/http-server/handlers"
	mw "todoApi/internal/http-server/middlewares"
	"todoApi/internal/storage"

	"github.com/gorilla/mux"
)

func New(storage *storage.Storage, log *slog.Logger) *mux.Router {
	router := mux.NewRouter()
	router.Use(mw.LogRequestsInfo(log))
	router.Handle("/task/create", mw.Post(handlers.CreateTask(log, storage)))
	router.Handle("/task", mw.Get(handlers.GetAllTasks(log, storage)))
	router.Handle("/task/{id}", mw.Get(handlers.GetTaskById(log, storage)))
	router.Handle("/task/delete/{id}", mw.Delete(handlers.DeleteTaskById(log, storage)))
	router.Handle("/task/set-completed/{id}", mw.Post(handlers.SetTaskCompletedById(log, storage)))

	return router
}
