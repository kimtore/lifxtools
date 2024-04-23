package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// NorthernLights generates "different" pixels (brighter, fuller, or different color)
// at random on a canvas filled with a base color.
type NorthernLights struct {
	Intensity float64 `json:"intensity"` // percentage chance to generate new pixel, from 0.0 (0%) to 1.0 (100%)
	Fade      float64 `json:"fade"`      // fade this many degrees each iteration
	Diversity float64 `json:"diversity"` // color hue diversity, in degrees maximum (usable range 0.0-180.0)
	Brighten  float64 `json:"brighten"`  // how much to brighten new dots, from 0.0 to 1.0.
	Saturate  float64 `json:"saturate"`  // how much to saturate new dots, from 0.0 to 1.0.
	Color     Color   `json:"color"`     // base color on the canvas
	degs      []float64
	dots      []colorful.Color
}

func init() {
	register("northernlights", func() Effect { return &NorthernLights{} })
}

func (e *NorthernLights) Init(pixels []colorful.Color) {
	fill(pixels, e.Color.Color)
	e.degs = make([]float64, len(pixels))
	e.dots = make([]colorful.Color, len(pixels))
}

func (e *NorthernLights) Draw(pixels []colorful.Color) {
	for i := range pixels {
		if e.degs[i] <= 0 && randomNonNegative() < e.Intensity {
			// Generate a new pixel within configured range
			h, c, l := e.Color.Hcl()
			h += random() * e.Diversity
			c += e.Saturate
			l += e.Brighten
			e.dots[i] = colorful.Hcl(h, c, l)
			e.degs[i] = 180.0
		}
		// Fade out a pixel that was generated earlier
		amplitude := math.Sin(e.degs[i] / Rad)
		pixels[i] = e.Color.BlendHcl(e.dots[i], amplitude)
		if e.degs[i] > 0 {
			e.degs[i] -= e.Fade
			if e.degs[i] < 0 {
				e.degs[i] = 0
			}
		}
	}
}
