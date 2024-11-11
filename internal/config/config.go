package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HttpServer `yaml:"http_server"`
	Postgres   `yaml:"postgres" env-required:"true"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	User     string `yaml:"user"`
	Password string `yaml:"password" `
	DBName   string `yaml:"db_name"`
	SslMode  string `yaml:"ssl_mode"`
}

func ParseConfig() *Config {
	var cfg Config

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error while cwd getting: ", err)
		return nil
	}

	configPath := filepath.Join(cwd, "./config/config.yaml")
	yamlConfig, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error while reading file: ", err)
		return nil
	}

	err = yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		log.Fatal("Error while parsing yaml: ")
		return nil
	}

	return &cfg
}
