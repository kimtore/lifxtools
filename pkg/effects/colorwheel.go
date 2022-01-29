package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

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

// Cycles through a list of colors indefinitely.
type ColorWheel struct {
	Hue       float64
	Chroma    float64
	Luminance float64
	Increment float64
}

func (e *ColorWheel) Draw(pixels []colorful.Color) {
	wheel := hclCircle(e.Hue, e.Chroma, e.Luminance, len(pixels))
	copy(pixels, wheel)
	e.Hue += e.Increment
}
