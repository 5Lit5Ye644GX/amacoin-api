package main

import (
	"fmt"

	"github.com/flibustier/multichain-client"
)

func init() {
	client := multichain.NewClient(
		"chain",
		"multichainrpc",
		"password",
		8080,
	).ViaNode(
		"localhost",
		8080,
	)

	obj, err := client.GetInfo() // Returns general information about this node and blockchain.
	if err != nil {
		fmt.Println("Il y a une erreur dans la connexion de la chaine. ")
		panic(err)
	}

	fmt.Println(obj)
}
