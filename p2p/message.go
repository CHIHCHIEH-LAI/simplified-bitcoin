package p2p

import "fmt"

type Message struct {
	Type    string `json:"type"`    // Type of the message (e.g. HEARTBEAT, TRANSACTION, BLOCK, etc)
	Sender  string `json:"sender"`  // Sender of the message
	Payload string `json:"payload"` // Payload of the message (as JSON string)
}

// NewMessage creates a new message
func NewMessage(messageType, sender, payload string) *Message {
	return &Message{
		Type:    messageType,
		Sender:  sender,
		Payload: payload,
	}
}

// Serialize serializes the message into a string
func (msg *Message) Serialize() string {
	return fmt.Sprintf("%s|%s|%s", msg.Type, msg.Sender, msg.Payload)
}

// DeserializeMessage deserializes the message from a string
func DeserializeMessage(data string) (Message, error) {
	var msg Message
	_, err := fmt.Sscanf(data, "%s|%s|%s", &msg.Type, &msg.Sender, &msg.Payload)
	if err != nil {
		return Message{}, fmt.Errorf("failed to deserialize message: %v", err)
	}

	return msg, nil
}
