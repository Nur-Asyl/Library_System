package readerRepo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"my_library/services/library/internal/domain/reader"
	"time"
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
	var reader reader.Reader
	err := r.db.QueryRowContext(ctx, "SELECT fio, address, nombil FROM readers WHERE nombil=$1", nombil).Scan(&reader.FIO, &reader.Address, &reader.NOMBIL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Reader not found")
		}
		return nil, err
	}

	return &reader, nil
}

func (r ReaderRepo) AcceptBook(ctx context.Context, name string, reader_nombil int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt := "UPDATE books SET nombil=NULL, date=NULL WHERE name=$1 AND nombil=$2"

	accepted, err := tx.PrepareContext(ctx, stmt)
	if err != nil {
		return err
	}

	res, err := accepted.ExecContext(ctx, name, reader_nombil)
	if err != nil {
		return err
	}

	if err = accepted.Close(); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("No such book or incorrect reader's nombil")
	}
	slog.Info("Records affected:", affected)

	return nil
}

func (r ReaderRepo) AssignBook(ctx context.Context, name string, reader_nombil int) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt := "UPDATE books SET nombil=$1, date=$2 WHERE name=$3 AND nombil IS NULL"

	assign, err := tx.PrepareContext(ctx, stmt)
	if err != nil {
		return err
	}

	res, err := assign.ExecContext(ctx, reader_nombil, time.Now(), name)
	if err != nil {
		return err
	}

	assign.Close()

	err = tx.Commit()
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("No such book or incorrect reader's nombil")
	}

	slog.Info("Records affected:", affected)

	return nil
}
