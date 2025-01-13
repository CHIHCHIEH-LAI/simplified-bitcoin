package network

import (
	"fmt"
	"log"
	"net"
)

// StartListener starts a TCP server and listens for incoming connections
func StartListener(port string) error {
	// Start listening on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start server on port %s: %v", port, err)
	}
	defer listener.Close()
	log.Printf("Server started on port %s\n", port)

	for {
		// Accept a single connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v\n", err)
			continue
		}

		// TODO: Handle the connection in a separate goroutine
		// Handle the connection (blocking call)
		handleConnection(conn)
	}
}

// handleConnection handles incoming connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read data from the connection
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Failed to read from connection: %v\n", err)
		return
	}

	// Print the received message
	message := string(buffer[:n])
	log.Printf("Received message: %s\n", message)
}
