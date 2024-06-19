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

func (r BookRepo) FindBooks(ctx context.Context, author, name string, year, invnom, nombil int) ([]book.Book, error) {
	var params []interface{}
	count := 1

	query := "SELECT author, name, year, date, invnom, nombil FROM books WHERE" + " "

	addAND := func() {
		if count > 1 {
			query += " AND "
		}
	}

	if author != "" {
		addAND()
		query += "author = $" + strconv.Itoa(count)
		count++
		params = append(params, author)
	}

	if name != "" {
		addAND()
		query += "name = $" + strconv.Itoa(count)
		count++
		params = append(params, name)
	}

	if year != 0 {
		addAND()
		query += "year = $" + strconv.Itoa(count)
		count++
		params = append(params, year)
	}

	if invnom != 0 {
		addAND()
		query += "invnom = $" + strconv.Itoa(count)
		count++
		params = append(params, invnom)
	}

	if nombil != 0 {
		addAND()
		query += "nombil = $" + strconv.Itoa(count)
		count++
		params = append(params, nombil)
	}

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Books not found")
		}
		return nil, err
	}
	defer rows.Close()

	var books []book.Book

	for rows.Next() {
		var newBook book.Book
		err := rows.Scan(&newBook.Author, &newBook.Name, &newBook.Year, &newBook.Date, &newBook.INVNOM, &newBook.NOMBIL)
		if err != nil {
			return nil, err
		}
		books = append(books, newBook)
	}

	if len(books) == 0 {
		return nil, errors.New("Books not found")
	}

	return books, nil
}
