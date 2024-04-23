package lifx

import (
	"encoding/binary"
	"io"
)

type GetVersionMessage struct {
	emptyMessage
}

func (m *GetVersionMessage) Type() uint16 {
	return MsgTypeGetVersion
}

type StateVersionMessage struct {
	emptyMessage
	Vendor    uint32
	Product   uint32
	Reserved6 [4]byte
}

func (m *StateVersionMessage) Unmarshal(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, m)
}

func (m *StateVersionMessage) Type() uint16 {
	return MsgTypeStateVersion
}

func (m *StateVersionMessage) Len() int {
	return binary.Size(m)
}
