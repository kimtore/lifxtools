package effects

import (
	"github.com/lucasb-eyer/go-colorful"
)

// Police generates "different" pixels (brighter, fuller, or different color)
// at random on a canvas filled with a base color.
type Police struct {
	Chroma    float64 `json:"chroma"`    // Saturation of colors, HCL space, from 0.0-1.0
	Luminance float64 `json:"luminance"` // Brightness of colors, HCL space, from 0.0-1.0
	pixelRoll []int
	index     int
	half      int
}

func init() {
	register("police", func() Effect { return &Police{} })
}

func (e *Police) Init(pixels []colorful.Color) {
	e.pixelRoll = []int{
		0, 0, -1, 0, -1, 0,
		0, 0, 1, 0, 1, 0,
	}
	e.half = len(pixels) / 2
}

func (e *Police) Draw(pixels []colorful.Color) {
	play := e.pixelRoll[e.index]
	black := colorful.Color{}
	red := colorful.Hcl(0, e.Chroma, e.Luminance)
	blue := colorful.Hcl(240, e.Chroma, e.Luminance)
	switch play {
	case 0:
		fill(pixels, black)
	case -1:
		fill(pixels[:e.half], red)
	case 1:
		fill(pixels[e.half:], blue)
	}
	e.index++
	e.index %= len(e.pixelRoll)
}
