package message

func NewJOINREQMessage(sender string) Message {
	return Message{
		Type:    "JOINREQ",
		Sender:  sender,
		Payload: "",
	}
}
