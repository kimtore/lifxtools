package lifx_test

import (
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMACAdressToFrameAddress(t *testing.T) {
	mac := []byte{0xd0, 0x73, 0x44, 0x55, 0xdf, 0x00}
	frameAddress := lifx.MACAdressToFrameAddress(mac)
	assert.Equal(t, uint64(0x00df554473d0), frameAddress)
}
