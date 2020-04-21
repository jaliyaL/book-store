package routers

import (
	"book-store/controller"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func Router() {
	fmt.Println("Hello world")

	r := mux.NewRouter()
	r.HandleFunc("/books", controller.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", controller.GetBook).Methods("GET")
	r.HandleFunc("/books", controller.AddBook).Methods("POST")
	r.HandleFunc("/books/{id}", controller.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", controller.RemoveBook).Methods("DELETE")

	r.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":2213", r))
}
