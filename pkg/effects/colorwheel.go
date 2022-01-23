package effects

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Generate an uniform set of colors across all hues using the HCL color space.
func HCLCircle(hue, chroma, luminance float64, steps int) []colorful.Color {
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
