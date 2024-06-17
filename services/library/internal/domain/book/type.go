package book

import (
	"errors"
	"time"
)

type Book struct {
	Author string
	Name   string
	Year   int
	INVNOM int
	Date   time.Time
	NOMBIL int
}

func NewBook(author, name string, year, invnom, nombil int) (*Book, error) {
	if err := checkAuthor(author); err != nil {
		return nil, err
	}
	if err := checkName(name); err != nil {
		return nil, err
	}

	return &Book{
		Author: author,
		Name:   name,
		Year:   year,
		INVNOM: invnom,
		NOMBIL: nombil,
	}, nil
}

//
//func (book Book) GetAuthor() string {
//	return book.author
//}
//
//func (book Book) GetName() string {
//	return book.name
//}
//
//func (book Book) GetYear() int {
//	return book.year
//}
//
//func (book Book) GetINVNOM() int {
//	return book.invnom
//}
//
//func (book Book) GetDate() time.Time {
//	return book.date
//}
//
//func (book Book) GetNOMBIL() int {
//	return book.nombil
//}

func checkAuthor(author string) error {
	if len(author) > 20 {
		return errors.New("Author exceeds 20 characters")
	}
	return nil
}

func checkName(name string) error {
	if len(name) > 40 {
		return errors.New("Name exceeds 40 characters")
	}
	return nil
}
