package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var frameAddressTests = []struct {
	payload      []byte
	frameAddress lifx.FrameAddress
}{
	{
		[]byte{
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		},
		lifx.FrameAddress{
			Target:      0,
			Reserved1:   [6]uint8{0, 0, 0, 0, 0, 0},
			Reserved2:   0,
			AckRequired: false,
			ResRequired: false,
			Sequence:    0,
		},
	},
	{
		[]byte{
			'\xaa', '\xbb', '\xcc', '\xdd', '\xee', '\xff', '\x00', '\x00',
			'\xaa', '\xbb', '\xcc', '\xdd', '\xee', '\xff', '\xaa', '\x71',
		},
		lifx.FrameAddress{
			Target:      0x0000ffeeddccbbaa,
			Reserved1:   [6]uint8{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			Reserved2:   0x2a,
			AckRequired: true,
			ResRequired: false,
			Sequence:    0x71,
		},
	},
}

// Test that a frame address header is correctly decoded into a FrameAddress struct.
func TestDecodeFrameAddress(t *testing.T) {
	for _, test := range frameAddressTests {
		a := assert.New(t)
		r := bytes.NewReader(test.payload)
		f, err := lifx.DecodeFrameAddress(r)
		a.Nil(err)
		a.Equal(test.frameAddress, *f)
	}
}

// Test that a frame address header is correctly encoded into a byte stream.
func TestEncodeFrameAddress(t *testing.T) {
	for _, test := range frameAddressTests {
		buf := make([]byte, 0, 16)
		w := bytes.NewBuffer(buf)
		a := assert.New(t)
		err := test.frameAddress.Write(w)
		a.Equal(16, w.Len())
		a.Nil(err)
		a.Equal(test.payload, w.Bytes())
	}
}
