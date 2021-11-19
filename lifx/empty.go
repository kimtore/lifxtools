package lifx

import (
	"encoding/binary"
	"io"
)

type emptyMessage struct {
}

func (m *emptyMessage) Unmarshal(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, m)
}

func (m *emptyMessage) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, m)
}

func (m *emptyMessage) Len() int {
	return binary.Size(m)
}
