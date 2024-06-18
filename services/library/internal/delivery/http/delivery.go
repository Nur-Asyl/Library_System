package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"my_library/services/library/configs"
	book2 "my_library/services/library/internal/domain/book"
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
	if err != nil {
		slog.Error("Error converting nombil to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reader, err := d.readerUC.FindReaderByNOMBIL(ctx, nombil)
	if err != nil {
		slog.Error("Error finding reader by nombil:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reader)
	if err != nil {
		slog.Error("Failed to encode searched reader:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	slog.Info("Reader searched", "status", http.StatusOK)
}

func (d *LibraryHTTPDelivery) AcceptBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	nombil, err := strconv.Atoi(r.URL.Query().Get("nombil"))
	if err != nil {
		slog.Error("Failed to conversate nombil to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.readerUC.AcceptBook(ctx, nombil); err != nil {
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

	book, err := d.GetBookParamsFromURl(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.readerUC.AssignBook(ctx, book.Author, book.Name, book.Year, book.INVNOM, book.NOMBIL); err != nil {
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

	book, err := d.GetBookParamsFromURl(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.bookUC.CreateBook(ctx, book.Author, book.Name, book.Year, book.INVNOM, book.NOMBIL); err != nil {
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

	book, err := d.GetBookParamsFromURl(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	books, err := d.bookUC.FindBooks(ctx, book.Author, book.Name, book.Year, book.INVNOM, book.NOMBIL)
	if err != nil {
		slog.Error("Failed to find books:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

	slog.Info("books searched", "status", http.StatusOK)
}

func (d *LibraryHTTPDelivery) Run(cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.Port)
	mux := http.NewServeMux()

	mux.HandleFunc("reader/create", d.CreateReaderHandler)
	mux.HandleFunc("reader/find", d.FindReaderByNOMBILHandler)
	mux.HandleFunc("reader/accept", d.AcceptBookHandler)
	mux.HandleFunc("reader/assign", d.AssignBookHandler)

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

func (d *LibraryHTTPDelivery) GetBookParamsFromURl(w http.ResponseWriter, r *http.Request) (*book2.Book, error) {
	var book *book2.Book
	author := r.URL.Query().Get("author")
	name := r.URL.Query().Get("name")
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		slog.Error("Failed to conversate year to int:", err)
		return nil, err
	}
	invnom, err := strconv.Atoi(r.URL.Query().Get("invnom"))
	if err != nil {
		slog.Error("Failed to conversate invnom to int:", err)
		return nil, err
	}
	nombil, err := strconv.Atoi(r.URL.Query().Get("nombil"))
	if err != nil {
		slog.Error("Failed to conversate nombil to int:", err)
		return nil, err
	}

	book.Author = author
	book.Name = name
	book.Year = year
	book.INVNOM = invnom
	book.NOMBIL = nombil
	return book, nil
}
