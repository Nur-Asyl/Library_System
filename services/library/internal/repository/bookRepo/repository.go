package bookRepo

import (
	"context"
	"database/sql"
	"errors"
	"my_library/services/library/internal/domain/book"
	"strconv"
)

type BookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (r BookRepo) CreateBook(ctx context.Context, book *book.Book) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO books (author, name, year, invnom) VALUES ($1, $2, $3, $4)", book.Author, book.Name, book.Year, book.INVNOM)
	if err != nil {
		return err
	}

	return nil
}

func (r BookRepo) FindBooks(ctx context.Context, b *book.Book) ([]*book.Book, error) {
	params := make([]any, 4)
	count := 1

	query := "SELECT author, name, year, invnom from books WHERE" + " "

	if b.Author != "" {
		query += "author = $" + strconv.Itoa(count) + " "
		count++
		params = append(params, b.Author)
	}

	if b.Name != "" {
		query += "name = $" + strconv.Itoa(count) + " "
		count++
		params = append(params, b.Name)
	}

	if b.Year != 0 {
		query += "year = $" + strconv.Itoa(count) + " "
		count++
		params = append(params, b.Year)
	}

	if b.INVNOM != 0 {
		query += "invnom = $" + strconv.Itoa(count) + " "
		count++
		params = append(params, b.INVNOM)
	}

	rows, err := r.db.QueryContext(ctx, query, params)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Books not found")
		}
		return nil, err
	}

	var books []*book.Book

	for rows.Next() {
		err := rows.Scan(&b.Author, &b.Name, &b.Year, &b.INVNOM, &b.Date, &b.NOMBIL)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return books, nil
}
