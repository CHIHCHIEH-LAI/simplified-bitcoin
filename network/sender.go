package network

import (
	"fmt"
	"log"
	"net"
)

// SendMessage sends a message data to the specified address
func SendMessageData(address string, messageData string) error {
	// Establish a connection to the remote address
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer conn.Close()

	// Send the message data
	_, err = conn.Write([]byte(messageData))
	if err != nil {
		return fmt.Errorf("failed to send message to %s: %v", address, err)
	}

	log.Printf("Message sent to %s: %s\n", address, messageData)
	return nil
}
