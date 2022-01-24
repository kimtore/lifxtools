package effects_test

import (
	"testing"

	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/stretchr/testify/assert"
)

func TestParseHCL(t *testing.T) {
	s := "hcl(0.15, 0.65, 0.99)"
	color, err := effects.ParseHCL(s)
	assert.NoError(t, err)

	h, c, l := color.Hcl()
	delta := 0.00001
	assert.InDelta(t, 0.15, h, delta)
	assert.InDelta(t, 0.65, c, delta)
	assert.InDelta(t, 0.99, l, delta)
}
