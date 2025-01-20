package membership

import (
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

// NewJOINREQMessage creates a new JOINREQ message
func NewJOINREQMessage(sender string) message.Message {
	return message.Message{
		Type:      message.JOINREQ,
		Sender:    sender,
		Payload:   "",
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// NewJOINRESPMessage creates a new JOINRESP message
func NewJOINRESPMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:      message.JOINRESP,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// NewHEARTBEATMessage creates a new HEARTBEAT message
func NewHEARTBEATMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:      message.HEARTBEAT,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}
