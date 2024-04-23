package virtbulb

import (
	"fmt"
	"sync"

	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

type canvas struct {
	pixels []lifx.HBSK
	lock   sync.Mutex
}

type Canvas interface {
	Draw(from, to int, hbsk lifx.HBSK)
	Pixels() []colorful.Color
}

func NewCanvas() Canvas {
	return &canvas{
		pixels: make([]lifx.HBSK, 0),
	}
}

func (c *canvas) Draw(from, to int, hbsk lifx.HBSK) {
	c.lock.Lock()
	defer c.lock.Unlock()
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

func lifxToColorful(hbsk lifx.HBSK) colorful.Color {
	return colorful.Hsv(
		(float64(hbsk.Hue)/65535.0)*360.0,
		float64(hbsk.Saturation)/65535.0,
		float64(hbsk.Brightness)/65535.0,
	)
}

func (c *canvas) Pixels() []colorful.Color {
	c.lock.Lock()
	defer c.lock.Unlock()
	colors := make([]colorful.Color, len(c.pixels))
	for i := range c.pixels {
		colors[i] = lifxToColorful(c.pixels[i])
	}
	return colors
}

func (c *canvas) Print() {
	c.lock.Lock()
	defer c.lock.Unlock()
	fmt.Println("---")
	for i := range c.pixels {
		color := lifxToColorful(c.pixels[i])
		s := textutil.SprintfHCL(color)
		fmt.Printf("%2d: %s\n", i, s)
	}
}
