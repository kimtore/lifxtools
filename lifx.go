package frame

import "io"

type frame struct {
	totsize     uint16 // 16 bits - Size of entire message in bytes (including this field)
	origin      uint8  // 2 bits - Message origin indicator (must be zero)
	tagged      bool   // 1 bit - Determines usage of the Frame address target field
	addressable bool   // 1 bit - Message includes a target address (must be one)
	protocol    uint16 // 12 bits - Protocol number (must be 1024 [decimal])
	source      uint32 // 32 bits - Source identifier: unique value set by the client, used by responses
	// frame address
	target      uint64 // 64 bits - 6 byte device address (MAC) (0 means all devices)
	fReserved   uint8  // 1 bits - reserved
	ackRequired bool   // 1 bits - Acknowledgment message required
	reqRequired bool   // 1 bits - Response message required
	sequence    uint8  // 8 bits - Wrap around message sequence number
	/* protocol frame */
	hReserved1 uint64 // 64 bits - Reserved
	mtype      uint16 // 16 bits - Message type determines the payload being used
	hReserved2 uint16 // 16 bits - Reserved
}

func (h *frame) Encode() {

}

func (h *frame) Decode() {

}

func NewDecoder(r io.Reader) *frame {
	return nil
}
