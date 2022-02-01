package virtbulb

import (
	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

type canvas struct {
	pixels []lifx.HBSK
}

type Canvas interface {
	Draw(from, to int, hbsk lifx.HBSK)
	Print()
}

func NewCanvas() Canvas {
	return &canvas{
		pixels: make([]lifx.HBSK, 0),
	}
}

func (c *canvas) Draw(from, to int, hbsk lifx.HBSK) {
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = len(c.pixels) - 1
	}
	if to > len(c.pixels)-1 {
		newpixels := make([]lifx.HBSK, to+1)
		copy(newpixels, c.pixels)
		c.pixels = newpixels
	}
	for ; from <= to; from++ {
		c.pixels[from] = hbsk
	}
}

func (c *canvas) Print() {
	for i := range c.pixels {
		color := colorful.Hsv(
			(float64(c.pixels[i].Hue)/65535.0)*360.0,
			float64(c.pixels[i].Saturation)/65535.0,
			float64(c.pixels[i].Brightness)/65535.0,
		)
		s := textutil.SprintfHCL(color)
		log.Debugf("%2d: %s", i, s)
	}
}
