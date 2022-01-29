package effects

import (
	"math/rand"

	"github.com/lucasb-eyer/go-colorful"
)

type NorthernLights struct {
	Intensity float64 `json:"intensity"` // chance to generate new pixel
	Fade      float64 `json:"fade"`      // how big percentage to fade out every iteration
	Cutoff    float64 `json:"cutoff"`    // how small difference is needed to zero out a dying color
	Diversity float64 `json:"diversity"` // color diversity, in degrees maximum
	Brighten  float64 `json:"brighten"`  // how much to brighten new dots
	Saturate  float64 `json:"saturate"`  // how much to saturate new dots
	Color     Color   `json:"color"`
}

func (e *NorthernLights) Init(pixels []colorful.Color) {
	for i := range pixels {
		pixels[i] = e.Color.Color
	}
}

func (e *NorthernLights) Draw(pixels []colorful.Color) {
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
