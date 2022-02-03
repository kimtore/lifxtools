package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Strobe shoots pixels across the canvas, fading them out in the process.
type Strobe struct {
	Color     Color   `json:"color"`     // Color of stroboscope
	Length    int     `json:"length"`    // Length of tail
	Diversity float64 `json:"diversity"` // Of color

	hue       float64
	direction int
	pos       int
	angle     []float64
	amplitude []float64
	dots      []colorful.Color
}

func init() {
	register("strobe", func() Effect { return &Strobe{} })
}

func (e *Strobe) Init(pixels []colorful.Color) {
	e.angle = make([]float64, len(pixels))
	e.amplitude = make([]float64, len(pixels))
	e.dots = make([]colorful.Color, len(pixels))
	e.hue, _, _ = e.Color.Hcl()
	e.direction = 1
	e.Draw(pixels)
}

func (e *Strobe) Draw(pixels []colorful.Color) {

	// Reference black can't be entirely black, since this will blend to red.
	h, c, l := e.Color.Hcl()
	black := colorful.Hcl(h, 0.01, 0)
	color := colorful.Hcl(e.hue, c, l)
	delta := 180.0 / float64(e.Length+1)

	for i := range pixels {
		if i == e.pos {
			e.angle[i] = 180 - e.angle[i]
			e.dots[i] = color
		}
		if e.angle[i] > 0 {
			// Fade out a pixel that was generated earlier
			amplitude := math.Sin(e.angle[i] / Rad)
			pixels[i] = black.BlendHcl(e.dots[i], amplitude)
			e.angle[i] -= delta
		} else {
			pixels[i] = black
		}
	}

	// Reverse direction if at either end
	e.pos += e.direction
	if e.pos >= len(pixels) {
		e.pos = len(pixels) - 1
	} else if e.pos < 0 {
		e.pos = 0
	} else {
		return
	}

	e.direction = -e.direction
	e.hue = h + (randomNonNegative() * e.Diversity * 180)
}
