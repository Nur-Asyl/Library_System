package main

import (
	"my_library/pkg/storage/postgres"
	"my_library/services/library/configs"
	"log"
	"log/slog"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration %+v", err)
	}

	storage, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Database %+v", err)
	}

	slog.Info("Successfully connected to database")

	readerRepo :=

}
