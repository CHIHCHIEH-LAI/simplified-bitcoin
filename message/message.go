package message

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

const (
	JOINREQ        = "JOINREQ"
	JOINRESP       = "JOINRESP"
	HEARTBEAT      = "HEARTBEAT"
	NEWTRANSACTION = "NEWTRANSACTION"
	NEWBLOCK       = "NEWBLOCK"
)

type Message struct {
	Type      string `json:"type"`      // Type of the message (e.g. HEARTBEAT, NEWTRANSACTION, etc)
	Sender    string `json:"sender"`    // Sender of the message
	Payload   string `json:"payload"`   // Payload of the message (as JSON string)
	Timestamp int64  `json:"timestamp"` // Timestamp of the message
}

// NewMessage creates a new message
func NewMessage(msgType, sender, payload string) *Message {
	return &Message{
		Type:      msgType,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// Serialize serializes the message into a string
func (msg *Message) Serialize() string {
	return fmt.Sprintf("%s|%s|%s|%d", msg.Type, msg.Sender, msg.Payload, msg.Timestamp)
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
