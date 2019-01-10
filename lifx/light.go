package lifx

// From https://lan.developer.lifx.com/docs/light-messages

import (
	"encoding/binary"
	"io"
)

// HSBK is used to represent the color and color temperature of a light.
//
// The color is represented as an HSB (Hue, Saturation, Brightness) value.
//
// The color temperature is represented in K (Kelvin) and is used to adjust the warmness / coolness of a white light, which is most obvious when saturation is close zero.
//
// Hue: range 0 to 65535
// Saturation: range 0 to 65535
// Brightness: range 0 to 65535
// Kelvin: range 2500° (warm) to 9000° (cool)
type HBSK struct {
	Hue        uint16
	Saturation uint16
	Brightness uint16
	Kelvin     uint16
}

// Sent by a client to change the light state.
//
// The duration is the color transition time in milliseconds.
//
// If the Frame Address res_required field is set to one (1) then the device will transmit a State message.
type SetColorMessage struct {
	Reserved uint8
	Color    HBSK
	Duration uint32
}

func DecodeSetColorMessage(r io.Reader) (*SetColorMessage, error) {
	m := &SetColorMessage{}
	err := binary.Read(r, binary.LittleEndian, m)
	return m, err
}

func (m *SetColorMessage) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, m)
}

func (m *SetColorMessage) Size() int {
	return 13
}
