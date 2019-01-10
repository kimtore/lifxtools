package lifx

import "io"

// Header combines Frame, FrameAddress and ProtocolHeader.
type Header struct {
	Frame          Frame
	FrameAddress   FrameAddress
	ProtocolHeader ProtocolHeader
}

// DecodeHeader reads an entire LIFX header and returns the resulting struct.
func DecodeHeader(r io.Reader) (*Header, error) {
	var err error

	f, err := DecodeFrame(r)
	if err != nil {
		return nil, err
	}

	fa, err := DecodeFrameAddress(r)
	if err != nil {
		return nil, err
	}

	ph, err := DecodeProtocolHeader(r)
	if err != nil {
		return nil, err
	}

	h := &Header{
		Frame:          *f,
		FrameAddress:   *fa,
		ProtocolHeader: *ph,
	}
	return h, nil
}

func (h *Header) Len() int {
	return h.Frame.Len() + h.FrameAddress.Len() + h.ProtocolHeader.Len()
}

func (h *Header) Write(w io.Writer) error {
	var err error

	err = h.Frame.Write(w)
	if err != nil {
		return err
	}
	err = h.FrameAddress.Write(w)
	if err != nil {
		return err
	}
	err = h.ProtocolHeader.Write(w)
	if err != nil {
		return err
	}

	return nil
}
