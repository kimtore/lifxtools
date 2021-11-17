package lifx

import "io"

type Payload interface {
	Len() int
	Write(io.Writer) error
	Type() uint16
}

// Packet combines Header and a payload to form a complete LIFX message.
type Packet struct {
	Header  Header
	Payload Payload
}

// DecodeHeader reads an entire LIFX message and returns the resulting struct.
func DecodePacket(r io.Reader) (*Packet, error) {
	var err error
	var payload Payload

	h, err := DecodeHeader(r)
	if err != nil {
		return nil, err
	}

	// Decode based on message type. Unsupported messages are read according to
	// size in the header field.
	switch h.ProtocolHeader.Type {
	case MsgTypeSetColorMessage:
		payload, err = DecodeSetColorMessage(r)
	case MsgTypeStateLabel, MsgTypeSetLabel:
		payload, err = DecodeSetLabelMessage(r)
	default:
		payloadSize := int(h.Frame.Size) - h.Len()
		if payloadSize > 0 {
			buf := make([]byte, 0, payloadSize)
			_, err = r.Read(buf)
		}
	}

	if err != nil {
		return nil, err
	}

	p := &Packet{
		Header:  *h,
		Payload: payload,
	}

	return p, nil
}

func (p *Packet) Len() int {
	return p.Header.Len() + p.Payload.Len()
}

func (p *Packet) Write(w io.Writer) error {
	var err error

	err = p.Header.Write(w)
	if err != nil {
		return err
	}

	err = p.Payload.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func (p *Packet) SetSize() {
	p.Header.Frame.Size = uint16(p.Len())
}

func NewPacket(payload Payload) *Packet {
	packet := &Packet{
		Header: Header{
			Frame: Frame{
				Addressable: true,
				Protocol:    ProtocolNumber,
			},
			ProtocolHeader: ProtocolHeader{
				Type: payload.Type(),
			},
		},
		Payload: payload,
	}
	packet.SetSize()
	return packet
}
