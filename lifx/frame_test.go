package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var frameTests = []struct {
	payload []byte
	frame   lifx.Frame
}{
	{
		[]byte{'\x31', '\x00', '\x00', '\x34', '\x00', '\x00', '\x00', '\x00'},
		lifx.Frame{49, 0, true, true, 1024, 0},
	},
	{
		[]byte{'\xff', '\x00', '\xaa', '\x7a', '\xd2', '\x02', '\x96', '\x49'},
		lifx.Frame{255, 1, true, true, 2730, 1234567890},
	},
	{
		[]byte{'\x11', '\xff', '\x30', '\xcf', '\x01', '\x23', '\x45', '\x67'},
		lifx.Frame{65297, 3, false, false, 3888, 1732584193},
	},
}

// Test that a frame header is correctly decoded into a Frame struct.
func TestDecodeFrame(t *testing.T) {
	for _, test := range frameTests {
		a := assert.New(t)
		r := bytes.NewReader(test.payload)
		f, err := lifx.DecodeFrame(r)
		a.Nil(err)
		a.Equal(test.frame, *f)
	}
}
