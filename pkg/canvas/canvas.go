package canvas

import (
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

type canvas struct {
	client lifx.Client
	size   int
	pixels []colorful.Color
	hbsk   []lifx.HBSK
}

type Canvas interface {
	Clear()
	Draw()
	Pixels() []colorful.Color
}

func New(client lifx.Client, size int) Canvas {
	return &canvas{
		client: client,
		size:   size,
		pixels: make([]colorful.Color, size),
		hbsk:   make([]lifx.HBSK, size),
	}
}

func tocolor(color colorful.Color) lifx.HBSK {
	const max = float64(65535)
	h, s, l := color.Hsl()
	return lifx.HBSK{
		Hue:        uint16(h * max),
		Saturation: uint16(s * max),
		Brightness: uint16(l * max),
	}
}

func (c *canvas) Clear() {
	for i := range c.pixels {
		c.pixels[i] = colorful.Color{}
		c.hbsk[i] = lifx.HBSK{}
	}
	c.client.SetColorZones(lifx.HBSK{}, 0, uint8(c.size-1), 0)
}

func (c *canvas) Draw() {
	for i, pixel := range c.pixels {
		hbsk := tocolor(pixel)
		if c.hbsk[i] == hbsk {
			continue
		}
		c.client.SetColorZones(hbsk, uint8(i), uint8(i), 0)
	}
}

func (c *canvas) Pixels() []colorful.Color {
	return c.pixels
}
