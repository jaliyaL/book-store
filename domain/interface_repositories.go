package domain

type BookRepo interface {
	GetBooks(bookId int) (res Book, err error)
}
