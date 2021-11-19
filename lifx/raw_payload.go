package lifx

import (
	"io"
)

// Catch-all for unknown payload types
type RawPayload struct {
	emptyMessage
	Data []byte
}

func DecodeRawPayload(r io.Reader, length int) (*RawPayload, error) {
	buf := make([]byte, length)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	return &RawPayload{
		Data: buf,
	}, nil
}

func (m *RawPayload) Write(w io.Writer) error {
	_, err := w.Write(m.Data)
	return err
}

func (m *RawPayload) Len() int {
	return len(m.Data)
}

func (m *RawPayload) Type() uint16 {
	return MsgTypeUnknown
}
