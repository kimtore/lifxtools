package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var setColorMessagePayload = []byte{
	'\x00', '\x55', '\x55', '\xFF', '\xFF', '\xFF', '\xFF', '\xAC', '\x0D',
	'\x00', '\x04', '\x00', '\x00',
}

// Test that a SetColor message payload is correctly decoded into a SetColorMessage struct.
func TestLightMessageDecoding(t *testing.T) {
	a := assert.New(t)
	r := bytes.NewReader(LightMessagePayload)
	msg, err := lifx.DecodeSetColorMessage(r)
	a.Nil(err)
	a.Equal(uint8(0), msg.Reserved)
	a.Equal(uint16(21845), msg.Color.Hue)
	a.Equal(uint16(65535), msg.Color.Saturation)
	a.Equal(uint16(65535), msg.Color.Brightness)
	a.Equal(uint16(3500), msg.Color.Kelvin)
}
