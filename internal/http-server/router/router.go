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
	router.Handle("/task/all", mw.Get(handlers.GetAllTasks(log, storage)))
	router.Handle("/task/one/{id}", mw.Get(handlers.GetTaskById(log, storage)))
	router.Handle("/task/delete/{id}", mw.Delete(handlers.DeleteTaskById(log, storage)))
	router.Handle("/task/set-completed/{id}", mw.Post(handlers.SetTaskCompletedById(log, storage)))
	router.Handle("/task/not-complited", mw.Get(handlers.GetUncomplitedTasks(log, storage)))
	router.Handle("/task/today", mw.Get(handlers.GetTodaysTasks(log, storage)))

	router.Handle("/auth/login", mw.Post(handlers.Login(log, storage)))
	router.Handle("/auth/registration", mw.Post(handlers.Registration(log, storage)))

	router.Handle("/user", mw.Get(handlers.GetAllUsers(log, storage)))
	router.Handle("/user/current", mw.AuthGuard(mw.Get(handlers.GetUserById(log, storage))))

	return router
}
