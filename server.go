package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/5Lit5Ye644GX/amacoin-api/blockchain"
	"github.com/5Lit5Ye644GX/amacoin-api/controllers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Controller must implements this functions
type Controller interface {
	GetAddresses(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
	GetStats(w http.ResponseWriter, r *http.Request)
	GetTransactions(w http.ResponseWriter, r *http.Request)
	CreateTransaction(w http.ResponseWriter, r *http.Request)
}

func main() {

	var controller Controller

	controller = controllers.Static{}

	port := flag.Int("p", 8080, "port to listen on")
	storage := flag.String("blockchain", "multichain", "blockchain is memory (no blockchain) / multichain / hyperledger")
	flag.Parse()

	if *storage == "multichain" {

		// Load .env file for setting Multichain parameters
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, please run the start script")
		}

		// Get Multichain RPC API Port
		mport, err := strconv.Atoi(os.Getenv("MULTICHAIN_PORT"))
		if err != nil || mport < 1 {
			log.Fatal("Cannot get multichain RPC port from .env")
		}

		// Instanciate our blockchain support
		blockchain := blockchain.NewBlockchain(os.Getenv("MULTICHAIN_RPC_PASSWORD"), mport)

		controller = controllers.Multichain{blockchain}
	}

	// Set routes
	router := mux.NewRouter()

	router.Methods("GET").Path("/addresses").HandlerFunc(controller.GetAddresses)
	router.Methods("GET").Path("/balance").HandlerFunc(controller.GetBalance)
	router.Methods("GET").Path("/stats").HandlerFunc(controller.GetStats)
	router.Methods("GET").Path("/transactions").HandlerFunc(controller.GetTransactions)
	router.Methods("POST").Path("/transactions").HandlerFunc(controller.CreateTransaction)

	host := fmt.Sprintf("http://localhost:%d", *port)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*", host},
		AllowedHeaders: []string{"authorization", "content-type"},
	})

	log.Printf("[OK] Server listening on %s\n", host)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), c.Handler(router)))
}
