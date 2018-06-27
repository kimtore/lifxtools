package lifx

import (
	"io"
	"io/ioutil"
)

type Frame struct {
	totsize     uint16 // 16 bits - Size of entire message in bytes (including this field)
	origin      uint8  // 2 bits - Message origin indicator (must be zero)
	tagged      bool   // 1 bit - Determines usage of the Frame address target field
	addressable bool   // 1 bit - Message includes a target address (must be one)
	protocol    uint16 // 12 bits - Protocol number (must be 1024 [decimal])
	source      uint32 // 32 bits - Source identifier: unique value set by the client, used by responses
	// Frame address
	target      uint64 // 64 bits - 6 byte device address (MAC) (0 means all devices)
	fReserved   uint8  // 1 bits - reserved
	ackRequired bool   // 1 bits - Acknowledgment message required
	reqRequired bool   // 1 bits - Response message required
	sequence    uint8  // 8 bits - Wrap around message sequence number
	/* protocol Frame */
	hReserved1 uint64 // 64 bits - Reserved
	mtype      uint16 // 16 bits - Message type determines the Payload being used
	hReserved2 uint16 // 16 bits - Reserved
}

func (h *Frame) Encode() {

}

func (h *Frame) Decode() (*Payload, error) {
	return nil, nil

}

// NewDecoder wraps an input stream, decodes the header,
// and returns a Payload-decoder that wraps the remaining input stream reader.
func NewDecoder(r io.Reader) *Frame {
	// Buffer the entire input before reading the fields
	data, err := ioutil.ReadAll(r)

	return nil
}

type Payload struct {
}

func (p *Payload) DecodePayload() (*Message, error) {
	if p == nil {
		return
	}

}

func (p *Payload) EncodePayload() (*Frame, error) {

}

func NewPayloadDecoder(r io.Reader) *Payload {
	return nil
}

type Message struct {
}
