package domain

type BookService interface {
	GetSelectedBook(bookId int) (res Book, err error)
}
