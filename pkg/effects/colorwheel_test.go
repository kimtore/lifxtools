package effects_test

import (
	"testing"

	"github.com/dorkowscy/lifxtool/pkg/effects"
)

func TestHCLCircle(t *testing.T) {
	colors := effects.HCLCircle(0, 0.3, 0.6, 200)
	for _, color := range colors {
		effects.PrintColorHSV(color)
	}
}
