package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"my_library/services/library/configs"
	"my_library/services/library/internal/use_case/book"
	"my_library/services/library/internal/use_case/reader"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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

func (d *LibraryHTTPDelivery) CreateReaderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	var requestData struct {
		FIO     string `json:"fio"`
		Address string `json:"address"`
		NOMBIL  int    `json:"nombil"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Error decoding request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.readerUC.CreateReader(ctx, requestData.FIO, requestData.Address, requestData.NOMBIL); err != nil {
		slog.Error("Error creating reader:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Reader created", "status", http.StatusCreated)
}

func (d *LibraryHTTPDelivery) FindReaderByNOMBILHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	nombil, err := strconv.Atoi(r.URL.Query().Get("nombil"))
	if err != nil || nombil == 0 {
		slog.Error("Error converting nombil to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("Converted nombil to int:", "nombil", nombil)

	reader, err := d.readerUC.FindReaderByNOMBIL(ctx, nombil)
	if err != nil {
		slog.Error("Error finding reader by nombil:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reader)
	if err != nil {
		slog.Error("Failed to encode searched reader:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Reader searched", "status", http.StatusOK)
}

func (d *LibraryHTTPDelivery) AcceptBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	var requestData struct {
		Name          string `json:"name"`
		Reader_nombil int    `json:"reader_nombil"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Error decoding request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Successfully decoded")

	if err := d.readerUC.AcceptBook(ctx, requestData.Name, requestData.Reader_nombil); err != nil {
		slog.Error("Failed to accept book:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Reader accepted book", "status", http.StatusOK)
}

func (d *LibraryHTTPDelivery) AssignBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	var requestData struct {
		Name          string `json:"name"`
		Reader_nombil int    `json:"reader_nombil"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Error decoding request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Successfully decoded")

	if err := d.readerUC.AssignBook(ctx, requestData.Name, requestData.Reader_nombil); err != nil {
		slog.Error("Failed to assign book to reader:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Reader assigned with book", "status", http.StatusOK)
}

func (d *LibraryHTTPDelivery) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	var requestData struct {
		Author string `json:"author"`
		Name   string `json:"name"`
		Year   int    `json:"year"`
		Invnom int    `json:"invnom"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Failed to decode request data:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("struct:", requestData)

	if err := d.bookUC.CreateBook(ctx, requestData.Author, requestData.Name, requestData.Year, requestData.Invnom); err != nil {
		slog.Error("Failed to create book:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Successfully created book")
}

func (d *LibraryHTTPDelivery) FindBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	var requestData struct {
		Author string `json:"author"`
		Name   string `json:"name"`
		Year   int    `json:"year"`
		INVNOM int    `json:"invnom"`
		NOMBIL int    `json:"nombil"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Error decoding request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Successfully decoded")

	books, err := d.bookUC.FindBooks(ctx, requestData.Author, requestData.Name, requestData.Year, requestData.INVNOM, requestData.NOMBIL)
	if err != nil {
		slog.Error("Failed to find books:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

	slog.Info("books searched", "status", http.StatusOK)
}

func (d *LibraryHTTPDelivery) Run(cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.Port)
	mux := http.NewServeMux()

	mux.HandleFunc("/reader/create", d.CreateReaderHandler)
	mux.HandleFunc("/reader/find", d.FindReaderByNOMBILHandler)
	mux.HandleFunc("/reader/accept", d.AcceptBookHandler)
	mux.HandleFunc("/reader/assign", d.AssignBookHandler)

	mux.HandleFunc("/book/create", d.CreateBookHandler)
	mux.HandleFunc("/book/find", d.FindBooks)

	go func() {
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			log.Fatalf("Failed to run http on port %s server %+v", cfg.Port, err)
		}
	}()
	slog.Info("Running http server", "port", cfg.Port)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT)
	<-quitCh

	slog.Info("Graceful shutdown")
}
