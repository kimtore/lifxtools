package lifx

import (
	"testing"
)

var LightMessage = []byte{
	'\x31', '\x00', '\x00', '\x34', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x66', '\x00', '\x00', '\x00', '\x00', '\x55', '\x55', '\xFF',
	'\xFF', '\xFF', '\xFF', '\xAC', '\x0D', '\x00', '\x04', '\x00', '\x00',
}

// assuming the server receives a message like LifxMessage above,
// the Decode function returns a Payload struct.
func TestDecode(t *testing.T) {
	// assrt := assert.New(t)

}
