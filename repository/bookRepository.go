package repository

import (
	"book-store/bootstrap"
	"book-store/domain"
	logger "github.com/sirupsen/logrus"
)

type BookRepo interface {
	GetBooks(bookId int) (res domain.Book, err error)
}

type BookRepoImplementation struct {
}

func (b BookRepoImplementation) GetBooks(bookId int) (res domain.Book, err error) {
	query := "select * from books where id = ?"

	rows := bootstrap.Conn.QueryRow(query, bookId)

	err = rows.Scan(&res.ID, &res.Title, &res.Author, &res.Year)
	if err != nil {
		logger.Error(err, err.Error())
	}

	return res, nil
}
