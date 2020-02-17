package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/book-store")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connection Established")
	}
	defer db.Close()

	fmt.Println("Hello world")

	r := mux.NewRouter()
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1234", r))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getBooks")
}

func getBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getBook")

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "You've requested book id: %s\n", id)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addBook")
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateBook")

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "You've requested book id: %s\n", id)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("removeBook")
}
