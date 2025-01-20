package message

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	JOINREQ        = "JOINREQ"
	JOINRESP       = "JOINRESP"
	HEARTBEAT      = "HEARTBEAT"
	NEWTRANSACTION = "NEWTRANSACTION"
)

type Message struct {
	Type      string `json:"type"`      // Type of the message (e.g. HEARTBEAT, TRANSACTION, BLOCK, etc)
	Sender    string `json:"sender"`    // Sender of the message
	Payload   string `json:"payload"`   // Payload of the message (as JSON string)
	Timestamp int64  `json:"timestamp"` // Timestamp of the message
}

// Serialize serializes the message into a string
func (msg *Message) Serialize() string {
	return fmt.Sprintf("%s|%s|%s,%d", msg.Type, msg.Sender, msg.Payload, msg.Timestamp)
}

// DeserializeMessage deserializes the message from a string
func DeserializeMessage(data string) (Message, error) {
	parts := strings.Split(data, "|")
	if len(parts) != 4 {
		return Message{}, fmt.Errorf("invalid message format")
	}

	timestamp, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return Message{}, fmt.Errorf("invalid timestamp format")
	}

	msg := Message{
		Type:      parts[0],
		Sender:    parts[1],
		Payload:   parts[2],
		Timestamp: timestamp,
	}

	return msg, nil
}
