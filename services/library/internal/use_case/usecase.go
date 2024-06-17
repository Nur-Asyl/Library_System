package use_case

import (
	"context"
	"my_library/services/library/internal/domain/book"
	"my_library/services/library/internal/domain/reader"
)

type BookUseCase interface {
	CreateBook(ctx context.Context, author, name string, year, invnom, nombil int) error
	FindBooks(ctx context.Context, author, name string, year, invnom, nombil int) (*[]book.Book, error)
}

type ReaderUseCase interface {
	CreateReader(ctx context.Context, fio, address string, nombil int) error
	FindReaderByNOMBIL(ctx context.Context, nombil int) (*reader.Reader, error)
	AcceptBook(ctx context.Context, nombil int) error
	AssignBook(ctx context.Context, author, name string, year, invnom, nombil int) error
}
