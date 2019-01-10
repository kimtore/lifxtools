package lifx

import (
	"encoding/binary"
	"io"
)

// The Frame Address section contains the following routing information:
//
// * Target device address
// * Acknowledgement message is required flag
// * State response message is required flag
// * Message sequence number
//
// The target device address is 8 bytes long, when using the 6 byte MAC address
// then left-justify the value and zero-fill the last two bytes. A target device
// address of all zeroes effectively addresses all devices on the local network.
// The Frame tagged field must be set accordingly.
//
// There are two flags that cause a LIFX device to send a message in response. In
// these cases, the source identifier in the response message will be set to the
// same value as that in the requesting message sent by the client.
//
// * _ack_required_ set to one (1) will cause the device to send an
//   Device::Acknowledgement message
// * _res_required_ set to one (1) within a Set message, e.g Light::SetPower will
//   cause the device to send the corresponding State message, e.g Light::StatePower
//
// The client can use acknowledgments to determine that the LIFX device has
// received a message. However, when using acknowledgments to ensure reliability
// in an over-burdened lossy network ... causing additional network packets may
// make the problem worse.
//
// Client that don't need to track the updated state of a LIFX device can choose
// not to request a response, which will reduce the network burden and may provide
// some performance advantage. In some cases, a device may choose to send a state
// update response independent of whether _res_required_ is set.
//
// The sequence number allows the client to provide a unique value, which will be
// included by the LIFX device in any message that is sent in response to a
// message sent by the client. This allows the client to distinguish between
// different messages sent with the same source identifier in the Frame. See
// _ack_required_ and _res_required_ fields in the Frame Address.
type FrameAddress struct {
	Target      uint64   // 64 bits - 6 byte device address (MAC address) or zero (0) means all devices
	Reserved1   [6]uint8 // 48 bits - must all be zero
	Reserved2   uint8    // 6 bits - reserved
	AckRequired bool     // 1 bits - Acknowledgment message required
	ResRequired bool     // 1 bits - Response message required
	Sequence    uint8    // 8 bits - Wrap around message sequence number
}

// DecodeFrameAddress reads sixteen bytes from a Reader, and returns the second part of a LIFX header.
func DecodeFrameAddress(r io.Reader) (*FrameAddress, error) {
	var err error
	var b uint8
	f := &FrameAddress{}

	err = binary.Read(r, binary.LittleEndian, &f.Target)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.LittleEndian, &f.Reserved1)
	if err != nil {
		return nil, err
	}

	// Read a byte of data and distribute it into fields
	err = binary.Read(r, binary.LittleEndian, &b)
	if err != nil {
		return nil, err
	}
	f.Reserved2 = (b >> 2) & 0x3f
	f.AckRequired = b&0x2 != 0
	f.ResRequired = b&0x1 != 0

	err = binary.Read(r, binary.LittleEndian, &f.Sequence)

	return f, err
}

func (f *FrameAddress) Write(w io.Writer) error {
	var err error
	var b uint8

	err = binary.Write(w, binary.LittleEndian, f.Target)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, f.Reserved1)
	if err != nil {
		return err
	}

	// Pack fields into a byte of data
	b |= ((f.Reserved2 & 0x3f) << 2)
	b |= (booltouint8(f.AckRequired) & 0x2)
	b |= (booltouint8(f.ResRequired) & 0x1)

	err = binary.Write(w, binary.LittleEndian, b)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, f.Sequence)
}

func (f *FrameAddress) Len() int {
	return 16
}
