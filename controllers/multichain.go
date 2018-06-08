package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/5Lit5Ye644GX/amacoin-api/repository"
)

// Multichain is an implementation of Controller with multichain integration
type Multichain struct {
	R *repository.MCRepository
}

// GetAddresses returns all addresses currently used
func (m Multichain) GetAddresses(w http.ResponseWriter, r *http.Request) {
	var addresses [2]struct {
		Address string `json:"address"`
	}
	addresses[0].Address = "13nNUaNU1XHKbBvPNQXtFnbVbgbD3vfhf6LTts"
	addresses[1].Address = "1ZESFph9SyhaxLrL1va4Qjq7cKVbuTh3BXozVj"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

// GetBalance returns the balance of the given address
func (m Multichain) GetBalance(w http.ResponseWriter, r *http.Request) {

	// need Authorization header
	if len(r.Header["Authorization"]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Missing Authorization header")
		return
	}
	address := r.Header["Authorization"][0]

	var balance struct {
		Balance float64 `json:"balance"`
	}
	balance.Balance = m.R.FetchBalance(address)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

// GetStats returns statistics about the Blockchain
func (m Multichain) GetStats(w http.ResponseWriter, r *http.Request) {
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
func (m Multichain) GetTransactions(w http.ResponseWriter, r *http.Request) {

	// need Authorization header
	if len(r.Header["Authorization"]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Missing Authorization header")
		return
	}
	address := r.Header["Authorization"][0]

	transactions := m.R.FetchTransactions(address)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// PostTransactions obviously creates a transaction from the given address
func (m Multichain) PostTransactions(w http.ResponseWriter, r *http.Request) {

	// need Authorization header
	if len(r.Header["Authorization"]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Missing Authorization header")
		return
	}
	// need Authorization header (with address$privkey)
	authorization := strings.Split(r.Header["Authorization"][0], "$")

	if len(authorization) < 2 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Wrong header")
		return
	}

	// need body with "to" (address) and "amount"
	type Message struct {
		To     int64  `json:"to"`
		Amount string `json:"amount"`
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// create a transaction to blockchain

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("impecc")
}
