package services

import (
	"book-store/domain"
	logger "github.com/sirupsen/logrus"
	//"book-store/bootstrap"
	"book-store/repository"
)

type BookService interface {
	GetSelectedBook(bookId int) (res domain.Book, err error)
}

var BRepo repository.BookRepo

func init() {
	BRepo = new(repository.BookRepoImplementation)
}

type BookServiceImplementation struct {
	Bepo repository.BookRepo
}

func (b BookServiceImplementation) GetSelectedBook(bookId int) (res domain.Book, err error) {
	b.Bepo = BRepo

	res, _ = b.Bepo.GetBooks(bookId)
	logger.Info("data", res)
	return res, nil
}

func Test(t string) string {
	return t
}
