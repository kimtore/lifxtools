package lifx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var LightMessage = []byte{
	'\x31', '\x00', '\x00', '\x34', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x66', '\x00', '\x00', '\x00', '\x00', '\x55', '\x55', '\xFF',
	'\xFF', '\xFF', '\xFF', '\xAC', '\x0D', '\x00', '\x04', '\x00', '\x00',
}

// assuming the server receives a message like LifxMessage above,
// the NewDecoder function returns a Frame struct.
func TestNewDecoder(t *testing.T) {
	assrt := assert.New(t)
	r := bytes.NewReader(LightMessage)
	nd := NewDecoder(r)
	assrt.NotNil(nd)
}
