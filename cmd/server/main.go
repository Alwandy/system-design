package main

import (
	"fmt"
	"github.com/Alwandy/system-design/api/v1/url"
	"github.com/Alwandy/system-design/pkg/dynamodb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/v1/url/newurl", url.NewUrlHandler).Methods("POST")
	r.HandleFunc("/api/v1/url/{id}", nil).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Create Tables if first time setup
	db.CreateTables()

	log.Printf("[INFO] Server started on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Julian system design")
}