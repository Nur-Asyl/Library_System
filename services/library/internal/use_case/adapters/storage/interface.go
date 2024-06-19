package storage

import (
	"context"
	"my_library/services/library/internal/domain/book"
	"my_library/services/library/internal/domain/reader"
)

type BookRepo interface {
	CreateBook(ctx context.Context, book *book.Book) error
	FindBooks(ctx context.Context, author, name string, year, invnom, nombil int) ([]book.Book, error)
}

type ReaderRepo interface {
	CreateReader(ctx context.Context, reader *reader.Reader) error
	FindReaderByNOMBIL(ctx context.Context, nombil int) (*reader.Reader, error)
	AcceptBook(ctx context.Context, name string, reader_nombil int) error
	AssignBook(ctx context.Context, name string, reader_nombil int) error
}
