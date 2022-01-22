package lifx

import (
	"bytes"
	"io"
)

type Payload interface {
	Len() int
	Type() uint16
	Write(io.Writer) error
	Unmarshal(io.Reader) error
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

	p := &Packet{}

	// Decode based on message type. Unsupported messages are read according to
	// size in the header field.
	switch h.ProtocolHeader.Type {
	case MsgTypeSetColorMessage:
		payload, err = DecodeSetColorMessage(r)
	case MsgTypeStateLabel, MsgTypeSetLabel:
		payload, err = DecodeSetLabelMessage(r)
	case MsgTypeSetAccessPoint:
		payload, err = DecodeSetAccessPointMessage(r)
	case MsgTypeStateVersion:
		payload = &StateVersionMessage{}
		err = payload.Unmarshal(r)
	default:
		payloadSize := int(h.Frame.Size) - h.Len()
		payload, err = DecodeRawPayload(r, payloadSize)
	}

	if err != nil {
		return nil, err
	}

	p.Header = *h
	p.Payload = payload

	return p, nil
}

func (p *Packet) Len() int {
	return p.Header.Len() + p.Payload.Len()
}

func (p *Packet) Write(w io.Writer) error {
	var err error
	buf := &bytes.Buffer{}

	err = p.Header.Write(buf)
	if err != nil {
		return err
	}

	err = p.Payload.Write(buf)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buf)
	return err
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
