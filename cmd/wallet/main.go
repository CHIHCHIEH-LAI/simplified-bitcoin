package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/wallet"
)

var (
	action     string // Action to perform: createWallet
	walletFile string // Filename for saving the wallet
)

func init() {
	// Define command-line flags
	flag.StringVar(&action, "action", "create", "Action to perform: 'createWallet'")
	flag.StringVar(&walletFile, "wallet", "wallet.json", "Filename for saving the wallet")
}

func main() {
	// Parse flags
	flag.Parse()

	switch action {
	case "createWallet":
		createWallet()
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
