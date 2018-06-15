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

//Stats used to return Blockchain's data
type Stats struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

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
func (mcr MCRepository) FetchInformations() []Stats {
	stats := make([]Stats, 2) // Stats that will store chain's data
	obj, err := mcr.m.GetInfo()
	if err != nil {
		log.Println("[ERROR] Fail to get information from multichain")
	}

	result := obj["result"].(map[string]interface{})
	stats[0].Key = "Name of the chain"
	stats[0].Value = result["chainname"].(string)
	stats[1].Key = "Blockchain height"
	stats[1].Value = fmt.Sprintf("#%.f", result["blocks"].(float64))
	return stats
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

//SendMoney is a function that allows one to send assets to another address
func (mcr MCRepository) SendMoney(from string, to string, amount float64, privkey string) error {

	msg := mcr.m.Command(
		"validateaddress",
		[]interface{}{
			privkey,
		},
	)

	resp, err := mcr.m.Post(msg)
	if err != nil || resp == nil {
		log.Printf("[ERROR] Could not validate address\n")
		return fmt.Errorf("can't validate the address")
	}

	a := resp.Result().(map[string]interface{})
	if a["isvalid"].(bool) != true {
		log.Printf("[ERROR] The address is not validated: %v\n", a)
		return fmt.Errorf("can't validate the address 3")
	}

	address := a["address"].(string)
	if address != from {
		log.Printf("[ERROR] Trying to send money without validated private key 4\n")
		return fmt.Errorf("can't validate the address 4")
	}

	_, err = mcr.m.ImportAddress(from, "", true)
	if err != nil {
		log.Printf("[ERROR] Cannot import address 5\n")
		return fmt.Errorf("cannot import address 5")
	}

	assets := make(map[string]float64, 1)
	assets[coinName] = amount

	blob, err := mcr.m.CreateRawSendFrom(from, to, assets)
	if err != nil {
		log.Printf("[ERROR] Cannot create raw transaction 6\n")
		return fmt.Errorf("cannot create Roz transaction 6")
	}

	grosblob, err := mcr.m.SignRawTransaction(blob.Result().(string), []*multichain.TxData{}, privkey)
	if err != nil {
		log.Printf("[ERROR] Cannot Sign raw transaction 7\n")
		return fmt.Errorf("cannot Sign Roz transaction 7")
	}

	_, err = mcr.m.SendRawTransaction(grosblob.Result().(map[string]interface{})["hex"].(string))
	if err != nil {
		log.Printf("[ERROR] SendRawTransaction error with address: %s private:%s\n%v", from, privkey, grosblob.Result())
		return fmt.Errorf("cannot Send Roz transaction 8")
	}

	return nil // Everything is all right.
}

//FetchAdresses is the function that returns the list of the addresses that are allowed to interact with the chain
func (mcr MCRepository) FetchAdresses() []string {
	tabret := make([]string, 0)
	params := []interface{}{"receive"}

	msg := mcr.m.Command( // It will do the manual command
		"listpermissions", // listpermissions that returns the allowed to receive a transaction
		params,            // Basically all the addresses of the network
	)
	coucou, erre := mcr.m.Post(msg)
	if erre != nil {
		log.Printf("[ERROR] Cannot execute listpermissions \n")
		return nil
	}

	for j := range coucou.Result().([]interface{}) { // Here we want to extract the addresses
		plop := coucou.Result().([]interface{})[j].(map[string]interface{})
		plip := plop["address"].(string)
		tabret = append(tabret, plip) // Adding the addresses
	}
	return tabret
}
