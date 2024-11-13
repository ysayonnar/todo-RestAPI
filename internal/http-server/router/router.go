package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"todoApi/internal/http-server/middlewares"
	"todoApi/internal/storage"
)

func tempHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		fmt.Fprint(w, "Hello world! id: ", v["id"])
	}
}

func New(storage *storage.Storage, log *slog.Logger) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.LogRequestsInfo(log))
	router.HandleFunc("/hello/{id}", tempHandler(log))

	return router
}
