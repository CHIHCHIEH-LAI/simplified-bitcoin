package network

import (
	"fmt"
	"log"
	"net"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
)

type Transmitter struct {
	MessageChannel <-chan message.Message // Channel to transmit messages
}

// NewTransmitter creates a new Transmitter instance
func NewTransmitter(messageChannel <-chan message.Message) *Transmitter {
	return &Transmitter{
		MessageChannel: messageChannel,
	}
}

// Run runs the transmitter
func (t *Transmitter) Run() {
	for msg := range t.MessageChannel {
		t.sendMessage(msg)
	}
}

// sendMessage sends a message to the specified address
func (t *Transmitter) sendMessage(msg message.Message) {
	// Establish a connection to the remote address
	conn, err := t.establishConnection(msg.Receipient)
	if err != nil {
		log.Printf("failed to establish connection: %v\n", err)
		conn.Close()
		return
	}

	// Serialize the message data
	msgData, err := msg.Serialize()
	if err != nil {
		log.Printf("failed to serialize message: %v\n", err)
		conn.Close()
		return
	}

	// Send the message data
	err = t.sendMessageData(conn, msgData)
	if err != nil {
		log.Printf("failed to send message: %v\n", err)
		conn.Close()
		return
	}

	conn.Close()
}

// establishConnection establishes a connection to the remote address
func (t *Transmitter) establishConnection(address string) (net.Conn, error) {
	// Establish a connection to the remote address
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	return conn, nil
}

// sendMessageData sends a message data to the specified connection
func (t *Transmitter) sendMessageData(conn net.Conn, messageData string) error {
	// Send the message data
	_, err := conn.Write([]byte(messageData))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent: %s\n", messageData)
	return nil
}

// Close closes the transmitter
func (t *Transmitter) Close() {
}
