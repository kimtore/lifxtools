package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var headerTests = []struct {
	payload []byte
	header  lifx.Header
}{
	{
		[]byte{
			'\x31', '\x00', '\x00', '\x34', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x66', '\x00', '\x00', '\x00', '\x00', '\x55', '\x55', '\xFF',
			'\xFF', '\xFF', '\xFF', '\xAC', '\x0D', '\x00', '\x04', '\x00', '\x00',
		},
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
	},
}

// Test that a header is correctly decoded into a Header struct.
func TestDecodeHeader(t *testing.T) {
	for _, test := range headerTests {
		a := assert.New(t)
		r := bytes.NewReader(test.payload)
		f, err := lifx.DecodeHeader(r)
		a.Nil(err)
		a.Equal(test.header, *f)
	}
}

// Test that a frame address header is correctly encoded into a byte stream.
func TestEncodeHeader(t *testing.T) {
	for _, test := range headerTests {
		buf := make([]byte, 0, 36)
		w := bytes.NewBuffer(buf)
		a := assert.New(t)
		err := test.header.Write(w)
		a.Equal(36, w.Len())
		a.Nil(err)
		a.Equal(test.payload[:36], w.Bytes())
	}
}
