package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

// Flags
var (
	mode    string // Mode: "start" or "send"
	port    string // Port to start the server
	target  string // Target node address (for send mode)
	message string // Message to send (for send mode)
)

func init() {
	// Define command-line flags
	flag.StringVar(&mode, "mode", "", "Mode to run the program in: 'start' or 'send'")
	flag.StringVar(&port, "port", "8080", "Port for the node to listen on (used in 'start' mode)")
	flag.StringVar(&target, "target", "", "Target node address to send the message to (used in 'send' mode')")
	flag.StringVar(&message, "message", "", "Message to send to the target node (used in 'send' mode')")
}

func main() {
	// Parse the flags
	flag.Parse()

	// Validate mode
	if mode == "" {
		fmt.Println("Error: Mode is required. Use '-mode=start' or '-mode=send'.")
		flag.Usage()
		os.Exit(1)
	}

	// Create a channel for messages (used in 'start' mode)
	messageChannel := make(chan string, 100)

	// Run based on the mode
	switch mode {
	case "start":
		// Start the node server
		go func() {
			err := network.RunListener(port, messageChannel)
			if err != nil {
				log.Fatalf("Failed to start node: %v\n", err)
			}
		}()

		// Process messages
		log.Println("Processing messages...")
		for msg := range messageChannel {
			handleMessage(msg)
		}

	case "send":
		// Validate required flags for send mode
		if target == "" || message == "" {
			log.Println("Error: Both '-target' and '-message' are required for 'send' mode.")
			flag.Usage()
			os.Exit(1)
		}

		// Send the message to the target node
		err := network.SendMessageData(target, message)
		if err != nil {
			log.Printf("Failed to send message: %v\n", err)
			os.Exit(1)
		}

		log.Println("Message sent successfully!")

	default:
		log.Println("Error: Invalid mode. Use '-mode=start' or '-mode=send'.")
		flag.Usage()
		os.Exit(1)
	}
}

// handleMessage handles incoming messages from the message channel
func handleMessage(msg string) {
	log.Printf("Received message: %+v\n", msg)
}
