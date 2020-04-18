package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

//var db *sql.DB
var err error

//func dbConn() (db *sql.DB) {
//	//dbDriver := "mysql"
//	//dbUser := "root"
//	//dbPass := "example"
//	//dbName := "book-store"
//	//dbHost := "172.18.0.3:3306"
//	//db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
//	if err != nil {
//		panic(err.Error())
//	}
//	return db
//}

func main() {

	//db, err := sql.Open("mysql", "root:example@tcp(172.18.0.4:3306)/book-store")

	//err = db.Ping()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//
	//defer db.Close()

	fmt.Println("Hello world")

	r := mux.NewRouter()
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	r.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":2213", r))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getBooks")
}

func getBook(w http.ResponseWriter, r *http.Request) {
	//db := dbConn()
	db, err := sql.Open("mysql", "root:example@tcp(172.18.0.2:3306)/book-store")

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("getBook")

	vars := mux.Vars(r)
	id := vars["id"]

	log.Println("id", id)

	var book Book

	query := "select * from books where id = ?"

	rows := db.QueryRow(query, id)

	err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		log.Fatal(err, err.Error())
	}
	////fmt.Println(book)
	//fmt.Println("bi")
	//err = db.QueryRow("select * from books where id=?", id).Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	//fmt.Println("hi")
	//if err != nil {
	//	log.Fatal(err)
	//}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(book)
	return
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
