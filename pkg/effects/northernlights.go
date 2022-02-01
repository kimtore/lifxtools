package effects

import (
	"math/rand"

	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

type NorthernLights struct {
	Intensity float64 `json:"intensity"` // chance to generate new pixel
	Fade      float64 `json:"fade"`      // how big percentage to fade out every iteration
	Cutoff    float64 `json:"cutoff"`    // how small difference is needed to zero out a dying color
	Diversity float64 `json:"diversity"` // color diversity, in degrees maximum
	Brighten  float64 `json:"brighten"`  // how much to brighten new dots
	Saturate  float64 `json:"saturate"`  // how much to saturate new dots
	Color     Color   `json:"color"`
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
func (e *NorthernLights) DrawLegacy(pixels []colorful.Color) {
	for i := range pixels {
		if e.Intensity < rand.Float64() {
			// Fade out existing pixels towards the base color
			if pixels[i].DistanceLab(e.Color.Color) > e.Cutoff {
				pixels[i] = pixels[i].BlendHcl(e.Color.Color, e.Fade)
			} else {
				pixels[i] = e.Color.Color
			}
		} else {
			// Generate a new pixel within configured range
			h, c, l := e.Color.Hcl()
			h += rnd() * e.Diversity
			c += e.Saturate
			l += e.Brighten
			col := colorful.Hcl(h, c, l)
			pixels[i] = e.Color.BlendHcl(col, rand.Float64())
		}
	}
}

func rnd() float64 {
	return 1 - rand.Float64() + rand.Float64()
}
