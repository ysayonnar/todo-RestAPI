package main

import (
	"log"
	"todoApi/internal/config"
	"todoApi/internal/logger"
)

func main() {
	cfg := config.ParseConfig()
	log.Println("Config parsed")

	log := logger.New(cfg)
	log.Debug("logger configured")
}
