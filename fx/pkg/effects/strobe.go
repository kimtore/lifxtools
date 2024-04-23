package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Strobe shoots pixels across the canvas, fading them out in the process.
type Strobe struct {
	Color     Color   `json:"color"`      // Color of stroboscope
	DutyCycle float64 `json:"duty-cycle"` // How long to stay "on" at a time, from 0-1 in 0.1 increments
	iteration float64
}

func init() {
	register("strobe", func() Effect { return &Strobe{} })
}

var (
	black = colorful.LinearRgb(0, 0, 0)
)

const (
	increase = 0.1
)

func (e *Strobe) Init(pixels []colorful.Color) {
	e.Draw(pixels)
}

func (e *Strobe) Draw(pixels []colorful.Color) {
	if e.iteration < e.DutyCycle {
		fill(pixels, e.Color.Color)
	} else {
		fill(pixels, black)
	}
	e.iteration = math.Mod(e.iteration+increase, 1)
}
