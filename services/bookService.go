package services

import (
	"book-store/domain"
	logger "github.com/sirupsen/logrus"
	//"book-store/bootstrap"
	repo "book-store/repository"
)

type BookServiceImplementation struct{}

func (b BookServiceImplementation) GetSelectedBook(bookId int) (res domain.Book, err error) {
	res, _ = repo.BookRepoImplementation{}.GetBooks(bookId)
	logger.Info("data", res)
	return res, nil
}

func Test(t string) string {
	return t
}
