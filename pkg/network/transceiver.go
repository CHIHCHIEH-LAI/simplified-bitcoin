package network

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"

type Transceiver struct {
	ReceiverChannel    chan *message.Message // Channel to receive messages
	TransmitterChannel chan *message.Message // Channel to transmit messages
	Receiver           *Receiver             // Receiver instance
	Transmitter        *Transmitter          // Transmitter instance
}

// NewTranceiver creates a new tranceiver
func NewTransceiver(port string) (*Transceiver, error) {
	tc := &Transceiver{
		ReceiverChannel:    make(chan *message.Message),
		TransmitterChannel: make(chan *message.Message),
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
func (tc *Transceiver) Run() {
	go tc.Receiver.Run()
	go tc.Transmitter.Run()
}

// Transmit sends a message
func (tc *Transceiver) Transmit(msg *message.Message) {
	tc.Transmitter.Transmit(msg)
}

// Receive receives a message
func (tc *Transceiver) Receive() (*message.Message, bool) {
	return tc.Receiver.Receive()
}

// Close closes the tranceiver
func (tc *Transceiver) Close() {
	// Close the receiver and transmitter
	tc.Receiver.Close()
	tc.Transmitter.Close()

	// Close the channels
	close(tc.ReceiverChannel)
	close(tc.TransmitterChannel)

}
