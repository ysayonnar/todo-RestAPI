package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"todoApi/internal/config"
	"todoApi/internal/http-server/router"
	"todoApi/internal/logger"
	"todoApi/internal/storage"
)



func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal("Error while parsing config: ", err)
	}
	log.Println("Config parsed")
	
	log := logger.New(cfg)
	
	log.Info("logger configured")

	storage, err := storage.NewStorageConnection(&cfg.Postgres)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}
	log.Info("database successfully connected")

	r := router.New(storage, log)
	server := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeOut,
	}

	log.Info("Server listens", slog.String("host", cfg.HttpServer.Address))
	if err = server.ListenAndServe(); err != nil {
		log.Error("Error while starting server", logger.Err(err))
	}
	log.Error("Server stopped!")
}
