package effects

import (
	"math/rand"

	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

// NorthernLights generates "different" pixels (brighter, fuller, or different color)
// at random on a canvas filled with a base color.
type NorthernLights struct {
	Intensity float64 `json:"intensity"` // percentage chance to generate new pixel, from 0.0 (0%) to 1.0 (100%)
	Fade      float64 `json:"fade"`      // percentage fadeout bright pixels every iteration
	Diversity float64 `json:"diversity"` // color hue diversity, in degrees maximum (usable range 0.0-180.0)
	Brighten  float64 `json:"brighten"`  // how much to brighten new dots, from 0.0 to 1.0.
	Saturate  float64 `json:"saturate"`  // how much to saturate new dots, from 0.0 to 1.0.
	Color     Color   `json:"color"`     // base color on the canvas
	fades     []float64
	dots      []colorful.Color
}

func init() {
	register("northernlights", func() Effect { return &NorthernLights{} })
}

func (e *NorthernLights) Init(pixels []colorful.Color) {
	fill(pixels, e.Color.Color)
	e.fades = make([]float64, len(pixels))
	e.dots = make([]colorful.Color, len(pixels))
}

func (e *NorthernLights) Draw(pixels []colorful.Color) {
	for i := range pixels {
		if e.fades[i] > 0 {
			// Fade out a pixel that was generated earlier
			pixels[i] = e.Color.BlendHcl(e.dots[i], e.fades[i])
			e.fades[i] -= e.Fade
		} else if rand.Float64() < e.Intensity {
			// Generate a new pixel within configured range
			h, c, l := e.Color.Hcl()
			h += rnd() * e.Diversity
			c += e.Saturate
			l += e.Brighten
			e.dots[i] = colorful.Hcl(h, c, l)
			e.fades[i] = 1
			log.Debugf("Generating new pixel at position %d: %s", i, textutil.SprintfHCL(e.dots[i]))
		} else {
			pixels[i] = e.Color.Color
		}
	}
}

func rnd() float64 {
	return 1 - rand.Float64() + rand.Float64()
}
