package message

import (
	"encoding/json"
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/utils"
)

const (
	JOINREQ        = "JOINREQ"
	JOINRESP       = "JOINRESP"
	HEARTBEAT      = "HEARTBEAT"
	NEWTRANSACTION = "NEWTRANSACTION"
	NEWBLOCK       = "NEWBLOCK"
	BLOCKCHAIN     = "BLOCKCHAIN"
)

type Message struct {
	Type       string `json:"type"`       // Type of the message (e.g. HEARTBEAT, NEWTRANSACTION, etc)
	Sender     string `json:"sender"`     // Sender of the message
	Receipient string `json:"receipient"` // Receipient of the message
	Payload    string `json:"payload"`    // Payload of the message (as JSON string)
	Timestamp  int64  `json:"timestamp"`  // Timestamp of the message
}

// NewMessage creates a new message
func NewMessage(msgType, sender, receipient, payload string) *Message {
	return &Message{
		Type:       msgType,
		Sender:     sender,
		Receipient: receipient,
		Payload:    payload,
		Timestamp:  utils.GetCurrentTimeInUnix(),
	}
}

// Serialize converts the Message to a JSON string
func (msg *Message) Serialize() (string, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %v", err)
	}
	return string(data), nil
}

// DeserializeMessage converts a JSON string back to a Message
func DeserializeMessage(data string) (*Message, error) {
	var msg Message
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %v", err)
	}
	return &msg, nil
}
