package lifx

import (
	"fmt"
	"io"
	"io/ioutil"
)

type bitStreamDecoder struct {
	dataStream []byte
	read       int // number of bits read
}

// Decode decodes the bitstream, and reads n bits from the underlying byte array
func (b *bitStreamDecoder) Decode(n int) []byte {
	nByte := read / 8 // The byte number in which we are currently positioned
	bitPos := read % 8
	data = byteStream[nByte]     // Get the byte which we are going to read from
	outByte := make([]byte, n/8) // Should be null bytes
	for i := range n {
		outByte[i/8] |= data & 2 << i
	}
	// TODO - The bits need to reversed if using little endian
	b.read += n
	return outByte
}

type Frame struct {
	reader      io.Reader // Wrapped input-stream
	totsize     uint16    // 16 bits - Size of entire message in bytes (including this field)
	origin      uint8     // 2 bits - Message origin indicator (must be zero)
	tagged      bool      // 1 bit - Determines usage of the Frame address target field
	addressable bool      // 1 bit - Message includes a target address (must be one)
	protocol    uint16    // 12 bits - Protocol number (must be 1024 [decimal])
	source      uint32    // 32 bits - Source identifier: unique value set by the client, used by responses
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

func (f *Frame) Encode() {

}

func (f *Frame) Decode() (*Payload, error) {
	// Buffer the entire input stream
	data, err := ioutil.ReadAll(h.reader)
	if err != nil {
		return nil, fmt.Errorf("Frame: Decode: error reading input stream: %v", err)
	}
	bd := bitStreamDecoder{data: data, read: 0}
	// Read the total Frame size
	f.totsize = bd.Decode(16) // Total frame size
	// Individually extract the bitfields
	f.origin = bd.Decode(2)
	f.tagged = bd.Decode(1)
	f.addressable = bd.Decode(1)
	f.protocol = bd.Decode(12)
	f.source = bd.Decode(32)
	// Frame address
	f.target = bd.Decode(64)
	f.fReserved = bd.Decode(1)
	f.ackRequired = bd.Decode(1)
	f.reqRequired = bd.Decode(1)
	f.sequence = bd.Decode(8)
	// protocol frame
	f.hReserved1 = bd.Decode(64)
	f.mtype = bd.Decode(16)
	f.hReserved2 = bd.Decode(16)
	return nil, nil

}

// NewDecoder wraps an input stream, decodes the header,
// and returns a Payload-decoder that wraps the remaining input stream reader.
func NewDecoder(r io.Reader) *Frame {
	return &Frame{reader: r}
}

type Payload struct {
}

func (p *Payload) DecodePayload() (*Message, error) {
	if p == nil {
		return nil, nil
	}
	return nil, nil
}

func (p *Payload) EncodePayload() (*Frame, error) {
	return nil, nil
}

func NewPayloadDecoder(r io.Reader) *Payload {
	return nil
}

type Message struct {
}
