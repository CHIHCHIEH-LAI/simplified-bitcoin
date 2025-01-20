package message

import (
	"fmt"
	"strings"
)

const (
	JOINREQ              = "JOINREQ"
	JOINRESP             = "JOINRESP"
	HEARTBEAT            = "HEARTBEAT"
	NEWTRANSACTION       = "NEWTRANSACTION"
	TRANSACTIONBROADCAST = "TRANSACTIONBROADCAST"
)

type Message struct {
	Type    string `json:"type"`    // Type of the message (e.g. HEARTBEAT, TRANSACTION, BLOCK, etc)
	Sender  string `json:"sender"`  // Sender of the message
	Payload string `json:"payload"` // Payload of the message (as JSON string)
}

// Serialize serializes the message into a string
func (msg *Message) Serialize() string {
	return fmt.Sprintf("%s|%s|%s", msg.Type, msg.Sender, msg.Payload)
}

// DeserializeMessage deserializes the message from a string
func DeserializeMessage(data string) (Message, error) {
	parts := strings.Split(data, "|")
	if len(parts) != 3 {
		return Message{}, fmt.Errorf("invalid message format")
	}

	msg := Message{
		Type:    parts[0],
		Sender:  parts[1],
		Payload: parts[2],
	}

	return msg, nil
}
