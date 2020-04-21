package services

import (
	"book-store/binding"
	"book-store/domain"
	logger "github.com/sirupsen/logrus"
)

type BookServiceImplementation struct {
	Bepo domain.BookRepo
}

func (b BookServiceImplementation) GetSelectedBook(bookId int) (res domain.Book, err error) {
	b.Bepo = binding.BRepo
	res, _ = b.Bepo.GetBooks(bookId)
	logger.Info("data", res)
	return res, nil
}

func Test(t string) string {
	return t
}
