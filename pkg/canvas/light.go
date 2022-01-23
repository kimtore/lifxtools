package canvas

import (
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

type light struct {
	client lifx.Client
	color  colorful.Color
	pixels []colorful.Color
}

func NewLight(client lifx.Client) Canvas {
	return &light{
		client: client,
	}
}

func (c *light) Fill(color colorful.Color) {
	c.color = color
}

func (c *light) Draw(fadeTime time.Duration) {
	hbsk := HBSK(c.color)
	log.Debugf("SetColor(hue=%v, sat=%v, bri=%v, fade=%v)\n", hbsk.Hue, hbsk.Saturation, hbsk.Brightness, fadeTime)
	c.client.SetColor(hbsk, fadeTime)
}

func (c *light) Set(color []colorful.Color) {
	if len(color) == 0 {
		c.color = colorful.Color{}
	} else {
		c.color = color[0]
	}
}

func (c *light) Pixels() []colorful.Color {
	return nil
}
