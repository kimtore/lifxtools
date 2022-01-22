package lifx

import (
	"encoding/binary"
	"io"
)

type MultiZoneApplicationRequest uint8

const (
	MultiZoneNoApply   MultiZoneApplicationRequest = 0
	MultiZoneApply     MultiZoneApplicationRequest = 1
	MultiZoneApplyOnly MultiZoneApplicationRequest = 2
)

type SetColorZonesMessage struct {
	emptyMessage
	StartIndex uint8
	EndIndex   uint8
	Color      HBSK
	Duration   uint32
	Apply      MultiZoneApplicationRequest
}

func DecodeSetColorZonesMessage(r io.Reader) (*SetColorZonesMessage, error) {
	m := &SetColorZonesMessage{}
	err := binary.Read(r, binary.LittleEndian, m)
	return m, err
}

func (m *SetColorZonesMessage) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, m)
}

func (m *SetColorZonesMessage) Len() int {
	return binary.Size(m)
}

func (m *SetColorZonesMessage) Type() uint16 {
	return MsgTypeSetColorZonesMessage
}
