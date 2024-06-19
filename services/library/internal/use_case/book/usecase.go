package book

import (
	"context"
	"log/slog"
	"my_library/services/library/internal/domain/book"
	"my_library/services/library/internal/use_case/adapters/storage"
)

type BookUseCase struct {
	bookRepo storage.BookRepo
}

func NewBookUseCase(bookRepo storage.BookRepo) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (uc BookUseCase) CreateBook(ctx context.Context, author, name string, year, invnom int) error {
	newBook, err := book.NewBook(author, name, year, invnom)
	if err != nil {
		slog.Error("Failed to construct book instance for creation")
		return err
	}

	if err := uc.bookRepo.CreateBook(ctx, newBook); err != nil {
		slog.Error("Failed to create book in db")
		return err
	}

	slog.Info("Book created in db")
	return nil
}

func (uc BookUseCase) FindBooks(ctx context.Context, author, name string, year, invnom, nombil int) ([]book.Book, error) {
	_, err := book.NewBook(author, name, year, invnom)
	if err != nil {
		slog.Error("Failed to construct book instance for search")
		return nil, err
	}

	books, err := uc.bookRepo.FindBooks(ctx, author, name, year, invnom, nombil)
	if err != nil {
		slog.Error("Failed to find books in db")
		return nil, err
	}

	return books, nil
}
