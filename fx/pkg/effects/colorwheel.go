package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// ColorWheel cycles through a list of colors indefinitely.
type ColorWheel struct {
	Chroma    float64 `json:"chroma"`    // Saturation of colors, HCL space, from 0.0-1.0
	Luminance float64 `json:"luminance"` // Brightness of colors, HCL space, from 0.0-1.0
	Increment float64 `json:"increment"` // Color hue is incremented by this much every iteration
	hue       float64
}

func init() {
	register("colorwheel", func() Effect { return &ColorWheel{} })
}

func (e *ColorWheel) Init(pixels []colorful.Color) {
	e.Draw(pixels)
}

func (e *ColorWheel) Draw(pixels []colorful.Color) {
	wheel := hclCircle(e.hue, e.Chroma, e.Luminance, len(pixels))
	copy(pixels, wheel)
	e.hue = math.Mod(e.hue+e.Increment, 360.0)
}

// Generate an uniform set of colors across all hues using the HCL color space.
func hclCircle(hue, chroma, luminance float64, steps int) []colorful.Color {
	const circle = 360.0
	colors := make([]colorful.Color, 0, steps)
	incr := circle / float64(steps)
	for steps > 0 {
		colors = append(colors, colorful.Hcl(hue, chroma, luminance))
		steps--
		hue = math.Mod(hue+incr, circle)
	}
	return colors
}
