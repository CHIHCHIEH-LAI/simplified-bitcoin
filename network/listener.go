package network

import (
	"fmt"
	"log"
	"net"
)

// StartListener starts a TCP server and listens for incoming connections
func RunListener(port string, messageChannel chan<- string) error {
	// Start listening on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start server on port %s: %v", port, err)
	}
	defer listener.Close()

	for {
		// Accept a single connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v\n", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go HandleConnection(conn, messageChannel)
	}
}

// handleConnection handles incoming connections
func HandleConnection(conn net.Conn, messageChannel chan<- string) {
	defer conn.Close()

	// Read data from the connection
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Failed to read from connection: %v\n", err)
		return
	}

	// Deserialize the received message data
	data := string(buffer[:n])
	messageChannel <- data
	log.Printf("Message data passed to handler: %+v\n", data)
}
