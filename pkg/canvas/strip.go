package canvas

import (
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

// LIFX Z strip, without extended multizone support.
type strip struct {
	client lifx.Client
	size   int
	pixels []colorful.Color
	hbsk   []lifx.HBSK
	cached []lifx.HBSK
}

func NewStrip(client lifx.Client, size int) Canvas {
	return &strip{
		client: client,
		size:   size,
		pixels: make([]colorful.Color, size),
		hbsk:   make([]lifx.HBSK, size+1),
		cached: make([]lifx.HBSK, size),
	}
}

func (c *strip) setColorZones(color lifx.HBSK, start, end uint8, fadeTime time.Duration) {
	log.Debugf("SetColorZones(%v, %v, %v, %v..%v, %v)\n", color.Hue, color.Saturation, color.Brightness, start, end, fadeTime)
	c.client.SetColorZones(color, start, end, fadeTime)
}

func (c *strip) Fill(color colorful.Color) {
	for i := range c.pixels {
		c.pixels[i] = color
	}
}

// Recursively look for identical or unchanged pixels in order to reduce packet count.
// We eliminate some packets by grouping together color zones that should have the same color.
func (c *strip) drawRange(start, i int, fadeTime time.Duration) int {
	if i >= c.size {
		c.setColorZones(c.hbsk[start], uint8(start), uint8(c.size-1), fadeTime)
		return i
	}
	// This canvas has a framebuffer which is aware of the colors on the strip,
	// and does not re-send colors that hasn't changed since the last call.
	if c.hbsk[i] == c.cached[i] {
		if start != i {
			c.setColorZones(c.hbsk[start], uint8(start), uint8(i-1), fadeTime)
		}
		return i + 1
	}
	c.cached[i] = c.hbsk[i]
	if c.hbsk[i] == c.hbsk[i+1] {
		return c.drawRange(start, i+1, fadeTime)
	}
	c.setColorZones(c.hbsk[start], uint8(start), uint8(i), fadeTime)
	return i + 1
}

func (c *strip) Draw(fadeTime time.Duration) {
	var i int
	for i = range c.pixels {
		c.hbsk[i] = HBSK(c.pixels[i])
	}
	for i = 0; i < c.size; {
		i = c.drawRange(i, i, fadeTime)
	}
}

func (c *strip) Set(pixels []colorful.Color) {
	if len(pixels) != c.size {
		return
	}
	copy(c.pixels, pixels)
}

func (c *strip) Pixels() []colorful.Color {
	return make([]colorful.Color, c.size)
}
