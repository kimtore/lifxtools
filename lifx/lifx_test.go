package lifx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var LifxMessage []byte = { 
'\x31', '\x00', '\x00', '\x34', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x66', '\x00', '\x00', '\x00', '\x00', '\x55', '\x55', '\xFF', '\xFF', '\xFF', '\xFF', '\xAC', '\x0D', '\x00', '\x04', '\x00', '\x00'
}

func TestAThing(t *testing.T) {
	assrt := assert.New(t)
}
