package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

func TestMACAdressToFrameAddress(t *testing.T) {
	mac := []byte{0xd0, 0x73, 0x44, 0x55, 0xdf, 0x00}
	frameAddress := lifx.MACAdressToFrameAddress(mac)
	assert.Equal(t, uint64(0x00df554473d0), frameAddress)
}

func TestWriteString(t *testing.T) {
	s := "hello, world!"
	buf := &bytes.Buffer{}

	buf.Reset()
	_, err := lifx.WriteString(buf, s, 0)
	assert.NoError(t, err)
	assert.Equal(t, "", buf.String())
	assert.Equal(t, 0, buf.Len())

	buf.Reset()
	_, err = lifx.WriteString(buf, "æøå", 2)
	assert.NoError(t, err)
	assert.Equal(t, "\xc3\x00", buf.String())
	assert.Equal(t, 2, buf.Len())

	buf.Reset()
	_, err = lifx.WriteString(buf, s, 16)
	assert.NoError(t, err)
	assert.Equal(t, "hello, world!\x00\x00\x00", buf.String())
	assert.Equal(t, 16, buf.Len())

	buf.Reset()
	_, err = lifx.WriteString(buf, s, 5)
	assert.NoError(t, err)
	assert.Equal(t, "hell\x00", buf.String())
	assert.Equal(t, 5, buf.Len())
}
