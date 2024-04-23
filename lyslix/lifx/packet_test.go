package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var packetTests = []struct {
	payload []byte
	packet  lifx.Packet
}{
	{
		[]byte{
			'\x31', '\x00', '\x00', '\x34', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x66', '\x00', '\x00', '\x00', '\x00', '\x55', '\x55', '\xFF',
			'\xFF', '\xFF', '\xFF', '\xAC', '\x0D', '\x00', '\x04', '\x00', '\x00',
		},
		lifx.Packet{
			lifx.Header{
				lifx.Frame{
					Size:        49,
					Origin:      0,
					Tagged:      true,
					Addressable: true,
					Protocol:    1024,
					Source:      0,
				},
				lifx.FrameAddress{
					Target:      0,
					Reserved1:   [6]uint8{0, 0, 0, 0, 0, 0},
					Reserved2:   0,
					AckRequired: false,
					ResRequired: false,
					Sequence:    0,
				},
				lifx.ProtocolHeader{
					Reserved1: 0,
					Type:      0x66,
					Reserved2: 0,
				},
			},
			&lifx.SetColorMessage{
				Duration: 1024,
				Color: lifx.HBSK{
					Hue:        21845,
					Saturation: 65535,
					Brightness: 65535,
					Kelvin:     3500,
				},
			},
		},
	},
}

// Test that a header is correctly decoded into a Header struct.
func TestDecodePacket(t *testing.T) {
	for _, test := range packetTests {
		a := assert.New(t)
		r := bytes.NewReader(test.payload)
		p, err := lifx.DecodePacket(r)
		a.Nil(err)
		a.Equal(test.packet, *p)
		a.NotNil(p)
	}
}

// Test that a frame address header is correctly encoded into a byte stream.
func TestEncodePacket(t *testing.T) {
	for _, test := range packetTests {
		buf := make([]byte, 0, 49)
		w := bytes.NewBuffer(buf)
		a := assert.New(t)
		err := test.packet.Write(w)
		a.Equal(49, w.Len())
		a.Nil(err)
		a.Equal(test.payload, w.Bytes())
	}
}
