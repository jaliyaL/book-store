package binding

import (
	"book-store/domain"
	"book-store/repository"
	"book-store/services"
)

var (
	BRepo    domain.BookRepo
	BService domain.BookService
)

func init() {
	BRepo = new(repository.BookRepoImplementation)
	BService = new(services.BookServiceImplementation)

}
