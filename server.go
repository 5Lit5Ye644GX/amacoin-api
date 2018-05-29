package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var blockchain *Blockchain

func main() {

	blockchain = NewBlockchain("password", 8080)

	port := 8080

	// Set routes
	router := mux.NewRouter()

	router.Methods("GET").Path("/addresses").HandlerFunc(GetAddresses)
	router.Methods("GET").Path("/balance").HandlerFunc(GetBalance)
	router.Methods("GET").Path("/stats").HandlerFunc(GetStats)
	router.Methods("GET").Path("/transactions").HandlerFunc(GetTransactions)
	router.Methods("POST").Path("/transactions").HandlerFunc(CreateTransaction)

	host := fmt.Sprintf("http://localhost:%d", port)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*", host},
		AllowedHeaders: []string{"authorization", "content-type"},
	})

	log.Printf("[OK] Server listening on %s\n", host)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), c.Handler(router)))
}
