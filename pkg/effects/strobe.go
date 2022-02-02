package effects

import (
	"github.com/lucasb-eyer/go-colorful"
)

// Strobe shoots pixels across the canvas, fading them out in the process.
type Strobe struct {
	//Chroma    float64 `json:"chroma"`    // Saturation of colors, HCL space, from 0.0-1.0
	//Luminance float64 `json:"luminance"` // Brightness of colors, HCL space, from 0.0-1.0
	Color     Color   `json:"color"`
	Length    int     `json:"length"`
	Intensity float64 `json:"intensity"` // Chance to generate new pixel
	Fade      float64 `json:"fade"`      // Fade out per iteration
	remaining int
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
	e.Draw(pixels)
}

func (e *Strobe) Draw(pixels []colorful.Color) {
	h, _, _ := e.Color.Hcl()
	black := colorful.Hcl(h, 0.01, 0)
	for i := range pixels {
		if i == e.pos && e.remaining > 0 {
			// Strobe flash in progress, draw next pixel
			e.dots[i] = e.Color.Color
			e.fades[i] = 1
			pixels[i] = e.dots[i]
			e.remaining--
		} else if e.fades[i] > 0 {
			// Fade out a pixel that was generated earlier
			pixels[i] = black.BlendHcl(e.dots[i], e.fades[i])
			e.fades[i] -= 1.0 / float64(e.Length)
		} else if e.remaining == 0 && randomNonNegative() < e.Intensity {
			e.remaining = e.Length
		} else {
			pixels[i] = black
		}
	}
	e.pos = (e.pos + 1) % len(pixels)
}
