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
	schemaregistry "github.com/landoop/schema-registry"
	"github.com/linkedin/goavro/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

var err error

type BookDetails struct {
	Title  string `json:"title"`
	Author struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	} `json:"author"`
	Year int `json:"year"`
}

type AddBookEvent struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
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

	//r.HandleFunc("/consume", consume).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":2213", r))
}

func consume() {

	fmt.Println("consume")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{"localhost:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "book-topic"
	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value), msg.Headers)
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
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

	//decode request
	bookDetails := BookDetails{}
	err := json.NewDecoder(r.Body).Decode(&bookDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//prepare data for kafka
	addBookEvent := AddBookEvent{}
	addBookEvent.Title = bookDetails.Title
	addBookEvent.Author = bookDetails.Author.Firstname + bookDetails.Author.Lastname
	addBookEvent.Year = bookDetails.Year

	byt, err := json.Marshal(addBookEvent)
	logger.Info("marshal request", byt)

	//get schema from schemaregistry
	client, _ := schemaregistry.NewClient("localhost:8081")
	//schema, _ := client.GetSchemaByID(10)
	schema, _ := client.GetSchemaBySubject("newBookAdded", 1)
	logger.Info("schema: ", schema.Schema)

	//goavro - encode process by using schema
	codec, err := goavro.NewCodec(schema.Schema)
	if err != nil {
		fmt.Println(err)
	}

	// Convert textual Avro data (in Avro JSON format) to native Go form
	native, _, err := codec.NativeFromTextual(byt)
	if err != nil {
		fmt.Println(err)
	}

	// Convert native Go form to binary Avro data
	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}

	logger.Info("decode request", bookDetails)
	//produce to kafka

	//kafka client - setup relevant config info
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

	//err = json.Unmarshal(byt, &bookDetails)
	//logger.Info("unmarshal request", bookDetails)

	//err = json.NewEncoder(w).Encode(book)

	topic := "book-topic"
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: 1,
		Key:       sarama.StringEncoder("test-key"),
		Value:     sarama.ByteEncoder(binary),
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
