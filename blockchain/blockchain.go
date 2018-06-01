package blockchain

import (
	"fmt"
	"log"

	"github.com/flibustier/multichain-client"
)

const (
	chainName = "test"
	chainUser = "multichainrpc"
	chainHost = "localhost"
	coinName  = "amacoin"
)

// Blockchain is our client wrapper
type Blockchain struct {
	m *multichain.Client
}

// Transaction contains information for a transaction
type Transaction struct {
	Date   int64   `json:"date"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// NewBlockchain initialize your blockchain
func NewBlockchain(password string, port int) *Blockchain {
	b := new(Blockchain)
	b.m = multichain.NewClient(
		chainName,
		chainUser,
		password,
		port,
	).ViaNode(
		chainHost,
		port,
	)

	// Returns general information about this node and blockchain.
	obj, err := b.m.GetInfo()
	if err != nil {
		log.Fatal("Fail to connect to multichain")
	}

	fmt.Println(obj)

	return b
}

// GetInfo returns informations about current Blockchain
func (b *Blockchain) GetInfo() {
	obj, err := b.m.GetInfo()
	if err != nil {
		log.Println("[ERROR] Fail to get information from multichain")
	}

	result := obj["result"].(map[string]interface{})

	fmt.Println(result["chainname"])
	fmt.Println(result["blocks"])

	obj, err = b.m.GetInfo()
	if err != nil {
		log.Println("[ERROR] Fail to get information from multichain")
	}
}

// GetBalance returns the amount of address' funds
func (b *Blockchain) GetBalance(address string) float64 {
	obj, err := b.m.GetAddressBalances(address)
	if err != nil {
		log.Printf("[ERROR] Fail to get balance for %s from multichain\n", address)
		return 0.0
	}

	result := obj.Result().([]interface{})
	if len(result) == 0 {
		return 0.0
	}
	return result[0].(map[string]interface{})["qty"].(float64)
}

// GetTransactions returns a list of Trasaction for the address
func (b *Blockchain) GetTransactions(address string) []Transaction {

	transactions := make([]Transaction, 0)

	obj, err := b.m.ListAddressTransactions(address, 100, 0, false)
	if err != nil {
		log.Printf("Could not list transactions from %s \n", address)
	}

	for index, element := range obj.Result().([]interface{}) {
		log.Println(index, element)
	}

	transactions = append(transactions, Transaction{1526978053, "1ZESFph9SyhaxLrL1va4Qjq7cKVbuTh3BXozVj", "13nNUaNU1XHKbBvPNQXtFnbVbgbD3vfhf6LTts", 10.01})

	return transactions
}
