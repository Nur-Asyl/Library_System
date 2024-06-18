package readerRepo

import (
	"context"
	"database/sql"
	"errors"
	"my_library/services/library/internal/domain/book"
	"my_library/services/library/internal/domain/reader"
)

type ReaderRepo struct {
	db *sql.DB
}

func NewReaderRepo(db *sql.DB) *ReaderRepo {
	return &ReaderRepo{db: db}
}

func (r ReaderRepo) CreateReader(ctx context.Context, reader *reader.Reader) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO readers (fio, address, nombil) VALUES ($1, $2, $3)", reader.FIO, reader.Address, reader.NOMBIL)
	if err != nil {
		return err
	}

	return nil
}

func (r ReaderRepo) FindReaderByNOMBIL(ctx context.Context, nombil int) (*reader.Reader, error) {
	var reader *reader.Reader
	err := r.db.QueryRowContext(ctx, "SELECT fio, address, nombil FROM readers WHERE nombil=$1", nombil).Scan(&reader.FIO, &reader.Address, &reader.NOMBIL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Reader not found")
		}
		return nil, err
	}

	return reader, nil
}

func (r ReaderRepo) AcceptBook(ctx context.Context, nombil int) error {
	row := r.db.QueryRowContext(ctx, "SELECT author, name, year, invnom FROM books WHERE nombil=$1", nombil)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (r ReaderRepo) AssignBook(ctx context.Context, book *book.Book) error {
	row := r.db.QueryRowContext(ctx, "SELECT fio, address, nombil FROM readers WHERE nombil=$1", book.NOMBIL)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
