package main

import (
	"log"
	"log/slog"
	"my_library/pkg/storage/postgres"
	"my_library/services/library/configs"
	"my_library/services/library/internal/delivery/http"
	bookRepo2 "my_library/services/library/internal/repository/bookRepo"
	readerRepo2 "my_library/services/library/internal/repository/readerRepo"
	"my_library/services/library/internal/use_case/book"
	"my_library/services/library/internal/use_case/reader"
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

	readerRepo := readerRepo2.NewReaderRepo(storage.GetDB())
	bookRepo := bookRepo2.NewBookRepo(storage.GetDB())

	readerUC := reader.NewReaderUseCase(readerRepo)
	bookUC := book.NewBookUseCase(bookRepo)

	delivery := http.NewLibraryHTTPDelivery(readerUC, bookUC)
	delivery.Run(cfg)
}
