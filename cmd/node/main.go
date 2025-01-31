package main

import (
	"flag"
	"log"
	"os"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/node"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/wallet"
)

var (
	port              string // Port to run the server
	address           string // Node address (e.g., "127.0.0.1:8080")
	bootstrapNodeAddr string // Address of the bootstrap node to join the network
	walletFile        string // Filename for saving the wallet
)

func init() {
	// Define command-line flags
	flag.StringVar(&port, "port", "8080", "Port for the node to listen on")
	flag.StringVar(&address, "address", "", "IP address of the node (e.g., 127.0.0.1:8080)")
	flag.StringVar(&bootstrapNodeAddr, "bootstrap", "", "Address of the bootstrap node to join the network (Optional)")
	flag.StringVar(&walletFile, "wallet", "wallet.json", "Filename for saving the wallet")
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Start the node
	startNode()
}

// startNode starts a new node
func startNode() {
	// Validate the address
	if address == "" {
		log.Println("Error: The node address is required. Use -address to specify it.")
		flag.Usage()
		os.Exit(1)
	}

	// Load the wallet from file
	w, err := wallet.LoadFromFile(walletFile)
	if err != nil {
		log.Fatalf("Failed to load wallet: %v\n", err)
	}

	// Get the address from the wallet
	address := w.GetAddress()

	// Create a new P2P node
	node, err := node.NewNode(address, port, address)
	if err != nil {
		log.Fatalf("Failed to create node: %v\n", err)
	}
	defer node.Close()

	// Start the node
	log.Printf("Starting node at %s...\n", address)
	node.Run(bootstrapNodeAddr)

	// Keep the server running
	select {}
}
