package lifx_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/stretchr/testify/assert"
)

// Onboard packet test fixture from https://github.com/tserong/lifx-hacks/blob/master/onboard.py
func TestDecodeOnboardMessage(t *testing.T) {
	onboard_packet := "\x86\x00\x00\x34\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x31\x01\x00\x00\x02"
	onboard_packet += "ssid" + strings.Repeat("\x00", 28)
	onboard_packet += "pass" + strings.Repeat("\x00", 60)
	onboard_packet += strings.Repeat("\x05", 1)
	r := strings.NewReader(onboard_packet)

	packet, err := lifx.DecodePacket(r)
	assert.NoError(t, err)
	assert.Equal(t, lifx.MsgTypeSetAccessPoint, int(packet.Header.ProtocolHeader.Type))
	payload := packet.Payload.(*lifx.SetAccessPointMessage)
	assert.Equal(t, packet.Len(), int(packet.Header.Frame.Size))
	assert.Equal(t, "ssid", payload.SSID)
	assert.Equal(t, "pass", payload.PSK)
	assert.Equal(t, lifx.Security_WPA2_AES_PSK, payload.Security)

	buf := &bytes.Buffer{}
	err = packet.Write(buf)
	assert.NoError(t, err)
	assert.Equal(t, onboard_packet, buf.String())
}
