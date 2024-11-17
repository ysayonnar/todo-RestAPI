package router

import (
	"log/slog"
	"todoApi/internal/http-server/handlers"
	"todoApi/internal/http-server/middlewares"
	"todoApi/internal/storage"

	"github.com/gorilla/mux"
)

func New(storage *storage.Storage, log *slog.Logger) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.LogRequestsInfo(log))
	router.Handle("/task/create", handlers.CreateTask(log, storage))

	return router
}
