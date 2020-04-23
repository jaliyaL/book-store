package main

import (
	"book-store/bootstrap"
	"book-store/domain"
	"book-store/services"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

var err error

type customer struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`
}

//type BookService interface {
//	GetSelectedBook(bookId int) (res domain.Book, err error)
//}

//type BookServiceImplementation struct{}
//
//func (b BookServiceImplementation) GetSelectedBook(bookId int) (res domain.Book, err error) {
//	query := "select * from books where id = ?"
//
//	rows := bootstrap.Conn.QueryRow(query, bookId)
//
//	err = rows.Scan(&res.ID, &res.Title, &res.Author, &res.Year)
//	if err != nil {
//		log.Fatal(err, err.Error())
//	}
//
//	return res, nil
//}

func main() {

	bootstrap.InitOptionalDB()

	defer bootstrap.CloseOptionalDB()

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

	vars := mux.Vars(r)
	id := vars["id"]

	bookId, _ := strconv.Atoi(id)

	//res, _ := services.BookServiceImplementation{}.GetSelectedBook(bookId)

	res, _ := services.BookServiceImplementation{}.GetSelectedBook(bookId)

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

func addBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addBook")

	//produce to kafka

	//setup relevant config info
	config := sarama.NewConfig()
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	//config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	//addresses of available kafka brokers
	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	//customer_details
	customerData := customer{
		Name:   "kusal",
		Age:    34,
		Height: 2.44,
		Weight: 132.45,
	}

	byt, err := json.Marshal(customerData)

	topic := "book-topic"
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: 1,
		//Key:       sarama.StringEncoder("test-key"),
		Value: sarama.ByteEncoder(byt),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
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
