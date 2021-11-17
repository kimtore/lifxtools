package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

var setLabelMessagePayload = []byte{
	0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77,
	0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

// Test that a SetLabel message payload is correctly decoded into a SetLabelMessage struct.
func TestLabelMessageDecoding(t *testing.T) {
	a := assert.New(t)
	r := bytes.NewReader(setLabelMessagePayload)
	msg, err := lifx.DecodeSetLabelMessage(r)
	a.Nil(err)
	a.Equal("hello, world!", msg.Label)
}

func TestLabelMessageEncoding(t *testing.T) {
	buf := &bytes.Buffer{}
	payload := &lifx.SetLabelMessage{
		Label: "hello, world!",
	}
	err := payload.Write(buf)
	assert.NoError(t, err)
	assert.Equal(t, setLabelMessagePayload, buf.Bytes())
}
