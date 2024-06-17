package storage

import (
	"context"
	"my_library/services/library/internal/domain/book"
	"my_library/services/library/internal/domain/reader"
)

type BookRepo interface {
	CreateBook(ctx context.Context, book *book.Book) error
	FindBooks(ctx context.Context, book *book.Book) ([]*book.Book, error)
}

type ReaderRepo interface {
	CreateReader(ctx context.Context, reader *reader.Reader) error
	FindReaderByNOMBIL(ctx context.Context, nombil int) (*reader.Reader, error)
	AcceptBook(ctx context.Context, nombil int) error
	AssignBook(ctx context.Context, book *book.Book) error
}
