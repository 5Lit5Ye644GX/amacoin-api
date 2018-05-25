package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// GetAddresses returns all addresses currently used
func GetAddresses(w http.ResponseWriter, r *http.Request) {
	var addresses [2]struct {
		Address string `json:"address"`
	}
	addresses[0].Address = "13nNUaNU1XHKbBvPNQXtFnbVbgbD3vfhf6LTts"
	addresses[1].Address = "1ZESFph9SyhaxLrL1va4Qjq7cKVbuTh3BXozVj"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

// GetBalance returns the balance of the given address
func GetBalance(w http.ResponseWriter, r *http.Request) {

	// need Authorization header

	var balance struct {
		Balance float64 `json:"balance"`
	}
	balance.Balance = 10.01

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

// GetStats returns statistics about the Blockchain
func GetStats(w http.ResponseWriter, r *http.Request) {
	var stats [3]struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	stats[0].Key = "Connected peers"
	stats[0].Value = "5"
	stats[1].Key = "Blockchain height"
	stats[1].Value = "#42"
	stats[2].Key = "Amacoin issued"
	stats[2].Value = "230"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetTransactions returns transactions concerning the given address
func GetTransactions(w http.ResponseWriter, r *http.Request) {

	// need Authorization header

	var transactions [2]struct {
		Date   int64   `json:"date"`
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	transactions[0].Date = 1526978053
	transactions[0].From = "1ZESFph9SyhaxLrL1va4Qjq7cKVbuTh3BXozVj"
	transactions[0].To = "13nNUaNU1XHKbBvPNQXtFnbVbgbD3vfhf6LTts"
	transactions[0].Amount = 10.01

	transactions[1].Date = 1526978150
	transactions[1].From = "13nNUaNU1XHKbBvPNQXtFnbVbgbD3vfhf6LTts"
	transactions[1].To = "1ZESFph9SyhaxLrL1va4Qjq7cKVbuTh3BXozVj"
	transactions[1].Amount = 0.01

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// CreateTransaction obviously creates a transaction from the given address
func CreateTransaction(w http.ResponseWriter, r *http.Request) {

	// need Authorization header (with address$privkey)

	// need body with "to" (address) and "amount"

	// create a transaction to blockchain

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("impecc")
}

func main() {

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
