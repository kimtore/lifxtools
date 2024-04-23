package lifx_test

import (
	"bytes"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

func TestSetColorZonesMessage_Write(t *testing.T) {
	expected := []byte{0xa, 0xb, 0xc, 0x0, 0xd, 0x0, 0xe, 0x0, 0xf, 0x0, 0x1a, 0x0, 0x0, 0x0, 0x01}

	packet := &lifx.SetColorZonesMessage{
		StartIndex: 0xa,
		EndIndex:   0xb,
		Color: lifx.HBSK{
			Hue:        0xc,
			Saturation: 0xd,
			Brightness: 0xe,
			Kelvin:     0xf,
		},
		Duration: 0x1a,
		Apply:    lifx.MultiZoneApply,
	}
	w := &bytes.Buffer{}
	a := assert.New(t)
	err := packet.Write(w)
	a.NoError(err)
	a.Equal(expected, w.Bytes())
}
