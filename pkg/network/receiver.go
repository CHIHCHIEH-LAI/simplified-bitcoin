package network

import (
	"fmt"
	"log"
	"net"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
)

type Receiver struct {
	Listener       net.Listener          // Listener to accept incoming connections
	MessageChannel chan *message.Message // Channel to send received messages
}

// NewReceiver creates a new Receiver instance
func NewReceiver(port string, messageChannel chan *message.Message) (*Receiver, error) {
	// Start listening on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to start server on port %s: %v", port, err)
	}

	return &Receiver{
		Listener:       listener,
		MessageChannel: messageChannel,
	}, nil
}

// Start starts the receiver
func (r *Receiver) Run() {
	for {
		// Accept a single connection
		conn, err := r.Listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v\n", err)
			conn.Close()
			continue
		}

		// Handle the connection in a separate goroutine
		r.handleConnection(conn)

		// Close the connection
		conn.Close()
	}
}

// Receives a message from the receiver
func (r *Receiver) Receive() (*message.Message, bool) {
	select {
	case msg := <-r.MessageChannel:
		return msg, true
	default:
		return nil, false // No message available
	}
}

// handleConnection handles incoming connections
func (r *Receiver) handleConnection(conn net.Conn) {
	// Read data from the connection
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Failed to read from connection: %v\n", err)
		return
	}

	// Deserialize the message
	msgData := string(buffer[:n])
	msg, err := message.DeserializeMessage(msgData)
	if err != nil {
		log.Printf("Failed to deserialize message: %v\n", err)
		return
	}
	log.Printf("Received message: %s\n", msgData)

	// Send the message to the message channel
	r.MessageChannel <- msg
}

// Close closes the receiver
func (r *Receiver) Close() {
	// Close the listener
	r.Listener.Close()

	// Close the message channel
	close(r.MessageChannel)
}
