package canvas

import (
	"fmt"
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

type light struct {
	client lifx.Client
	color  colorful.Color
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
	hbsk := ToHBSK(c.color)
	c.client.SetColor(hbsk, fadeTime)
	fmt.Printf("SetColor(%v, %v, %v, %v)\n", hbsk.Hue, hbsk.Saturation, hbsk.Brightness, fadeTime)
}

func (c *light) Set(_ []colorful.Color) {

}

func (c *light) Pixels() []colorful.Color {
	return nil
}
