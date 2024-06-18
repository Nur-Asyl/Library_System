package http

import (
	"fmt"
	"log"
	"log/slog"
	"my_library/services/library/configs"
	"my_library/services/library/internal/use_case/book"
	"my_library/services/library/internal/use_case/reader"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type LibraryHTTPDelivery struct {
	readerUC *reader.ReaderUseCase
	bookUC   *book.BookUseCase
}

func NewLibraryHTTPDelivery(readerUC *reader.ReaderUseCase, bookUC *book.BookUseCase) *LibraryHTTPDelivery {
	return &LibraryHTTPDelivery{
		readerUC: readerUC,
		bookUC:   bookUC,
	}
}

func (d *LibraryHTTPDelivery) Run(cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.Port)
	mux := http.NewServeMux()

	go func() {
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			log.Fatalf("Failed to run http server %+v", err)
		}
	}()
	slog.Info("Running http server", "port", cfg.Port)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT)
	<-quitCh

	slog.Info("Graceful shutdown")
}
