package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Cycles through a list of colors indefinitely.
type ColorWheel struct {
	Hue       float64 `json:"hue"`
	Chroma    float64 `json:"chroma"`
	Luminance float64 `json:"luminance"`
	Increment float64 `json:"increment"`
}

func init() {
	register("colorwheel", func() Effect { return &ColorWheel{} })
}

func (e *ColorWheel) Init(pixels []colorful.Color) {
	e.Draw(pixels)
}

func (e *ColorWheel) Draw(pixels []colorful.Color) {
	wheel := hclCircle(e.Hue, e.Chroma, e.Luminance, len(pixels))
	copy(pixels, wheel)
	e.Hue = math.Mod(e.Hue+e.Increment, 360.0)
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
