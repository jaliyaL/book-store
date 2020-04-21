package controller

import (
	"book-store/domain"
	"book-store/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var BookSvr domain.BookService

func GetBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getBooks")
}

func GetBook(w http.ResponseWriter, r *http.Request) {

	//BookSvr := binding.BService

	vars := mux.Vars(r)
	id := vars["id"]

	bookId, _ := strconv.Atoi(id)

	res, _ := services.BookServiceImplementation{}.GetSelectedBook(bookId)

	//res, _ := BookSvr.GetSelectedBook(bookId)

	book := domain.Book{}
	book.Author = res.Author
	book.ID = res.ID
	book.Title = res.Title
	book.Year = res.Year

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)
	return
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addBook")
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateBook")

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "You've requested book id: %s\n", id)
}

func RemoveBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("removeBook")
}
