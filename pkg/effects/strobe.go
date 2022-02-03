package effects

import (
	"github.com/lucasb-eyer/go-colorful"
)

// Strobe shoots pixels across the canvas, fading them out in the process.
type Strobe struct {
	Color     Color   `json:"color"`     // Color of stroboscope
	Length    int     `json:"length"`    // Length of tail
	Intensity float64 `json:"intensity"` // Chance to generate new pixel
	Fade      float64 `json:"fade"`      // Fade out per iteration

	direction int
	pos       int
	fades     []float64
	dots      []colorful.Color
}

func init() {
	register("strobe", func() Effect { return &Strobe{} })
}

func (e *Strobe) Init(pixels []colorful.Color) {
	e.fades = make([]float64, len(pixels))
	e.dots = make([]colorful.Color, len(pixels))
	e.direction = 1
	e.Draw(pixels)
}

func (e *Strobe) Draw(pixels []colorful.Color) {
	h, _, _ := e.Color.Hcl()

	// Reference black can't be entirely black, since this will blend to red.
	black := colorful.Hcl(h, 0.01, 0)
	for i := range pixels {
		if i == e.pos {
			// Strobe flash in progress, draw next pixel
			e.dots[i] = e.Color.Color
			e.fades[i] = 1
			pixels[i] = e.dots[i]
		} else if e.fades[i] > 0 {
			// Fade out a pixel that was generated earlier
			pixels[i] = black.BlendHcl(e.dots[i], e.fades[i])
			e.fades[i] -= 1.0 / float64(e.Length)
		} else {
			pixels[i] = black
		}
	}

	e.pos += e.direction
	if e.pos >= len(pixels) {
		e.pos = len(pixels) - 1
		e.direction = -1
	} else if e.pos < 0 {
		e.pos = 0
		e.direction = 1
	}
}
