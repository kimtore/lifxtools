package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var protocolHeaderTests = []struct {
	payload        []byte
	protocolHeader lifx.ProtocolHeader
}{
	{
		[]byte{
			'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
			'\x00', '\x00', '\x00', '\x00',
		},
		lifx.ProtocolHeader{
			Reserved1: 0,
			Type:      0,
			Reserved2: 0,
		},
	},
	{
		[]byte{
			'\xaa', '\xbb', '\xcc', '\xdd', '\xee', '\xff', '\x00', '\x11',
			'\xaa', '\x00', '\xbb', '\xcc',
		},
		lifx.ProtocolHeader{
			Reserved1: 0x1100ffeeddccbbaa,
			Type:      0x00aa,
			Reserved2: 0xccbb,
		},
	},
}

// Test that a protocol header is correctly decoded into a ProtocolHeader struct.
func TestDecodeProtocolHeader(t *testing.T) {
	for _, test := range protocolHeaderTests {
		a := assert.New(t)
		r := bytes.NewReader(test.payload)
		f, err := lifx.DecodeProtocolHeader(r)
		a.Nil(err)
		a.Equal(test.protocolHeader, *f)
	}
}

// Test that a frame address header is correctly encoded into a byte stream.
func TestEncodeProtocolHeader(t *testing.T) {
	for _, test := range protocolHeaderTests {
		buf := make([]byte, 0, 12)
		w := bytes.NewBuffer(buf)
		a := assert.New(t)
		err := test.protocolHeader.Write(w)
		a.Equal(12, w.Len())
		a.Nil(err)
		a.Equal(test.payload, w.Bytes())
	}
}
