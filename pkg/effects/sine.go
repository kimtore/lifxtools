package effects

import (
	"math"

	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/lucasb-eyer/go-colorful"
)

// Cycles through a list of colors indefinitely.
type Sine struct {
	Color     Color   `json:"color"`
	Hue       float64 `json:"hue"`
	Chroma    float64 `json:"chroma"`
	Luminance float64 `json:"luminance"`
	Increment float64 `json:"increment"`
	deg       float64
}

const Rad = 360.0 / (math.Pi * 2)

func init() {
	register("sine", &Sine{})
}

func (e *Sine) Init(pixels []colorful.Color) {
	e.Draw(pixels)
}

func (e *Sine) Draw(pixels []colorful.Color) {
	amplitude := math.Sin(e.deg / Rad)
	h, c, l := e.Color.Hcl()
	h += amplitude * e.Hue
	c += amplitude * e.Chroma
	l += amplitude * e.Luminance
	color := colorful.Hcl(h, c, l)
	textutil.PrintColorHCL(color)
	fill(pixels, color)
	e.deg = math.Mod(e.deg+e.Increment, 360.0)
}
