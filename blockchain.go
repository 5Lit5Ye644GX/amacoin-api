package main

import (
	"fmt"
	"log"

	"github.com/flibustier/multichain-client"
)

const (
	chainName = "test"
	chainUser = "multichainrpc"
	chainHost = "localhost"
)

// Blockchain is our client wrapper
type Blockchain struct {
	m *multichain.Client
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
		fmt.Errorf("Fail to connect to multichain")
	}

	fmt.Println(obj)
}
