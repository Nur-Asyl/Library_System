package reader

import (
	"context"
	"log/slog"
	"my_library/services/library/internal/domain/book"
	"my_library/services/library/internal/domain/reader"
	"my_library/services/library/internal/use_case/adapters/storage"
)

type ReaderUseCase struct {
	readerRepo storage.ReaderRepo
}

func NewReaderUseCase(readerRepo storage.ReaderRepo) *ReaderUseCase {
	return &ReaderUseCase{readerRepo: readerRepo}
}

func (uc ReaderUseCase) CreateReader(ctx context.Context, fio, address string, nombil int) error {
	newReader, err := reader.NewReader(fio, address, nombil)
	if err != nil {
		slog.Error("Failed to construct reader")
		return err
	}

	if err := uc.readerRepo.CreateReader(ctx, newReader); err != nil {
		slog.Error("Failed to create reader in db")
		return err
	}

	return nil

}

func (uc ReaderUseCase) FindReaderByNOMBIL(ctx context.Context, nombil int) (*reader.Reader, error) {
	reader, err := uc.readerRepo.FindReaderByNOMBIL(ctx, nombil)
	if err != nil {
		slog.Error("Failed to find reader by nombil in db")
		return nil, err
	}

	return reader, nil
}

func (uc ReaderUseCase) AcceptBook(ctx context.Context, nombil int) error {
	if err := uc.readerRepo.AcceptBook(ctx, nombil); err != nil {
		slog.Error("Failed to accept bookRepo in db")
		return err
	}

	return nil
}

func (uc ReaderUseCase) AssignBook(ctx context.Context, author, name string, year, invnom, nombil int) error {
	book, err := book.NewBook(author, name, year, invnom, nombil)
	if err != nil {
		slog.Error("Failed to construct bookRepo")
		return err
	}

	if err := uc.readerRepo.AssignBook(ctx, book); err != nil {
		slog.Error("Failed to assignbook")
		return err
	}

	return nil
}
