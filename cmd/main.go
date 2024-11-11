package main

import (
	"log"
	"todoApi/internal/config"
)

func main() {
	cfg := config.ParseConfig()
	log.Println(cfg.Address)

	// TODO: setup logger
}
