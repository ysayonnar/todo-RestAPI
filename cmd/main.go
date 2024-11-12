package main

import (
	"log"
	"os"
	"todoApi/internal/config"
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
	_ = storage
}
