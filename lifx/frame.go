package lifx

import (
	"encoding/binary"
	"io"
)

// The Frame section contains information about the following:
//
// * Size of the entire message
// * LIFX Protocol number: must be 1024 (decimal)
// * Use of the Frame Address target field
// * Source identifier
//
// The tagged field is a boolean flag that indicates whether the Frame Address
// target field is being used to address an individual device or all devices. For
// discovery using Device::GetService the tagged field should be set to one (1)
// and the target should be all zeroes. In all other messages the tagged field
// should be set to zero (0) and the target field should contain the device MAC
// address. The device will then respond with a Device::StateService message,
// which will include its own MAC address in the target field. In all subsequent
// messages that the client sends to the device, the target field should be set to
// the device MAC address, and the tagged field should be set to zero (0).
//
// The source identifier allows each client to provide an unique value, which will
// be included by the LIFX device in any message that is sent in response to a
// message sent by the client. If the source identifier is a non-zero value, then
// the LIFX device will send a unicast message to the IP address and port of the
// client that sent the originating message. If the source identifier is a zero
// value, then the LIFX device may send a broadcast message that can be received
// by all clients on the same sub-net. See _ack_required_ and _res_required_
// fields in the Frame Address.
type Frame struct {
	Size        uint16 // 16 bits - Size of entire message in bytes (including this field)
	Origin      uint8  // 2 bits - Message origin indicator (must be zero)
	Tagged      bool   // 1 bit - Determines usage of the Frame address target field
	Addressable bool   // 1 bit - Message includes a target address (must be one)
	Protocol    uint16 // 12 bits - Protocol number (must be 1024 [decimal])
	Source      uint32 // 32 bits - Source identifier: unique value set by the client, used by responses
}

// DecodeFrame reads eight bytes from a Reader, and returns the first part of a LIFX header.
func DecodeFrame(r io.Reader) (*Frame, error) {
	var err error
	var w uint16
	f := &Frame{}

	err = binary.Read(r, binary.LittleEndian, &f.Size)
	if err != nil {
		return nil, err
	}

	// Read a word of data and distribute it into fields
	err = binary.Read(r, binary.LittleEndian, &w)
	if err != nil {
		return nil, err
	}
	f.Origin = uint8(w>>14) & 0x3
	f.Tagged = w&0x2000 != 0
	f.Addressable = w&0x1000 != 0
	f.Protocol = w & 0xfff

	err = binary.Read(r, binary.LittleEndian, &f.Source)

	return f, err
}

// EncodeFrame writes a frame to a Writer, totaling eight bytes.
func EncodeFrame(f *Frame, wr io.Writer) (error) {
	var err error
	var w uint16

	err = binary.Write(wr, binary.LittleEndian, f.Size)
	if err != nil {
		return err
	}

	// Pack fields into a word of data
	w |= (uint16(f.Origin & 0x3)) << 14
	w |= (booltouint16(f.Tagged) & 0x2000)
	w |= (booltouint16(f.Addressable) & 0x1000)
	w |= (f.Protocol & 0xfff)

	err = binary.Write(wr, binary.LittleEndian, w)
	if err != nil {
		return err
	}

	err = binary.Write(wr, binary.LittleEndian, f.Source)

	return nil
}
