package repository

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

// MCRepository is our client wrapper
type MCRepository struct {
	m *multichain.Client
}

// Transaction contains information for a transaction
type Transaction struct {
	Date   float64 `json:"date"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// NewMCRepository initialize the connection to multichain
func NewMCRepository(password string, port int) *MCRepository {
	mcr := new(MCRepository)
	mcr.m = multichain.NewClient(
		chainName,
		chainUser,
		password,
		port,
	).ViaNode(
		chainHost,
		port,
	)

	// Returns general information about this node and blockchain.
	obj, err := mcr.m.GetInfo()
	if err != nil {
		log.Fatal("Fail to connect to multichain")
	}

	fmt.Println(obj)

	return mcr
}

// FetchInformations returns informations about current Blockchain
func (mcr MCRepository) FetchInformations() {
	obj, err := mcr.m.GetInfo()
	if err != nil {
		log.Println("[ERROR] Fail to get information from multichain")
	}

	result := obj["result"].(map[string]interface{})

	fmt.Println(result["chainname"])
	fmt.Println(result["blocks"])

	obj, err = mcr.m.GetInfo()
	if err != nil {
		log.Println("[ERROR] Fail to get information from multichain")
	}
}

// FetchBalance returns the amount of address' funds
func (mcr MCRepository) FetchBalance(address string) float64 {
	obj, err := mcr.m.GetAddressBalances(address)
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

// FetchTransactions returns a list of Trasaction for the address
func (mcr MCRepository) FetchTransactions(address string) []Transaction {

	transactions := make([]Transaction, 0)

	obj, err := mcr.m.ListAddressTransactions(address, 100, 0, false)
	if err != nil || obj == nil {
		log.Printf("[ERROR] Could not list transactions from %s \n", address)
		return transactions
	}

	// ListAddressTransactions returns all transactions (including permissions)
	for _, element := range obj.Result().([]interface{}) {
		e := element.(map[string]interface{})
		// if addresses is empty, it's not a transaction we want
		if len(e["addresses"].([]interface{})) > 0 {
			balance := e["balance"].(map[string]interface{})["assets"].([]interface{})
			// if the balance is empty, it's not a transaction we want
			if len(balance) > 0 {
				amount := balance[0].(map[string]interface{})["qty"].(float64)
				from := e["addresses"].([]interface{})[0].(string)
				to := e["myaddresses"].([]interface{})[0].(string)
				if amount < 0 {
					swap := from
					from = to
					to = swap
					amount *= -1
				}
				t := Transaction{Date: e["time"].(float64), From: from, To: to, Amount: amount}
				transactions = append(transactions, t)
			}
		}
	}

	return transactions
}
