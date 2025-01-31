package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/wallet"
)

var (
	address           string  // Address of the node (e.g., "127.0.0.1:8080")
	bootstrapNodeAddr string  // Address of the bootstrap node to join the network
	action            string  // Action to perform: createWallet, createTx
	walletFile        string  // Filename for saving the wallet
	recipient         string  // Recipient address for the transaction
	amount            float64 // Amount to send in the transaction
	fee               float64 // Transaction fee
)

func init() {
	// Define command-line flags
	flag.StringVar(&address, "address", "", "IP address of the node (e.g., 127.0.0.1:8080")
	flag.StringVar(&bootstrapNodeAddr, "bootstrap", "", "Address of the bootstrap node to join the network")
	flag.StringVar(&action, "action", "create", "Action to perform: 'createWallet', 'createTx'")
	flag.StringVar(&walletFile, "wallet", "wallet.json", "Filename for saving the wallet")
	flag.StringVar(&recipient, "recipient", "", "Recipient address for the transaction")
	flag.Float64Var(&amount, "amount", 0.0, "Amount to send in the transaction")
	flag.Float64Var(&fee, "fee", 0.0, "Transaction fee")
}

func main() {
	// Parse flags
	flag.Parse()

	switch action {
	case "createWallet":
		createWallet()
	case "createTx":
		createTransaction()
	default:
		fmt.Println("Invalid action. Use 'createWallet'")
		flag.Usage()
		os.Exit(1)
	}
}

func createWallet() {
	// Create a new wallet
	w := wallet.NewWallet()
	fmt.Printf("New wallet created!\nAddress: %s\n", w.GetAddress())

	// Save the wallet to file
	err := w.SaveToFile(walletFile)
	if err != nil {
		log.Fatalf("Failed to save wallet: %v\n", err)
	}

	fmt.Printf("Wallet saved to '%s'\n", walletFile)
}

func createTransaction() {
	// Load the wallet from file
	w, err := wallet.LoadFromFile(walletFile)
	if err != nil {
		log.Fatalf("Failed to load wallet: %v\n", err)
	}

	// Create a new transaction
	tx, err := w.CreateTransaction(recipient, amount, fee)
	if err != nil {
		log.Fatalf("Failed to create transaction: %v\n", err)
	}
	fmt.Printf("Transaction created!\nID: %s\n", tx.TransactionID)

	// Send the transaction to the network
	err = w.SendTransaction(tx, address, bootstrapNodeAddr)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v\n", err)
	}
}
