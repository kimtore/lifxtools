package lifx

import (
	"encoding/binary"
	"io"
)

// The Protocol header contains the following information about the message:
//
// * Message type which determines what action to take (based on the Payload)
type ProtocolHeader struct {
	Reserved1 uint64 // 64 bits - Reserved
	Type      uint16 // 16 bits - Message type determines the Payload being used
	Reserved2 uint16 // 16 bits - Reserved
}

func DecodeProtocolHeader(r io.Reader) (*ProtocolHeader, error) {
	h := &ProtocolHeader{}
	err := binary.Read(r, binary.LittleEndian, h)
	return h, err
}

func (p *ProtocolHeader) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, p)
}

func (p *ProtocolHeader) Len() int {
	return 12
}
