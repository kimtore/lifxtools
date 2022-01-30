package effects_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/dorkowscy/lifxtool/pkg/effects"
)

func TestSine(t *testing.T) {
	for deg := 0.0; deg < 360; deg++ {
		n := math.Sin(deg / effects.Rad)
		fmt.Println(n)
	}
}
