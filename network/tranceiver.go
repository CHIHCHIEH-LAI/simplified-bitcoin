package network

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"

type Tranceiver struct {
	ReceiverChannel    chan message.Message // Channel to receive messages
	TransmitterChannel chan message.Message // Channel to transmit messages
	Receiver           *Receiver            // Receiver instance
	Transmitter        *Transmitter         // Transmitter instance
}

// NewTranceiver creates a new tranceiver
func NewTranceiver(port string) (*Tranceiver, error) {
	tc := &Tranceiver{
		ReceiverChannel:    make(chan message.Message),
		TransmitterChannel: make(chan message.Message),
	}

	var err error
	tc.Receiver, err = NewReceiver(port, tc.ReceiverChannel)
	if err != nil {
		return nil, err
	}
	tc.Transmitter = NewTransmitter(tc.TransmitterChannel)

	return tc, nil
}

// Run runs the tranceiver
func (tc *Tranceiver) Run() {
	go tc.Receiver.Run()
	go tc.Transmitter.Run()
}

// Close closes the tranceiver
func (tc *Tranceiver) Close() {
	close(tc.ReceiverChannel)
	close(tc.TransmitterChannel)
	tc.Receiver.Close()
	tc.Transmitter.Close()
}
